package mc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/spf13/cast"
	"io/ioutil"
	"reflect"
	"strings"
)

//从json文件读数据到struct
func JsonFileUnmarshal(file string, obj interface{}) error {
	if jdData, err := ioutil.ReadFile(file); err != nil{
		return err
	}else{
		jdData = bytes.TrimPrefix(jdData, []byte("\xef\xbb\xbf"))
		return json.Unmarshal(jdData, obj)
	}
}
//从json文件读数据到struct(有特殊处理)
func McJsonFileUnmarshal(file string, obj interface{}) error {
	if jdData, err := ioutil.ReadFile(file); err != nil{
		return err
	}else{
		jdData = bytes.TrimPrefix(jdData, []byte("\xef\xbb\xbf"))
		var data map[string]interface{}
		if err := json.Unmarshal(jdData, &data); err != nil{
			return err
		}
		return reflectSetStruct(obj, data)
	}
}


func reflectSetStruct(obj interface{},data map[string]interface{}) (err error){
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	sfName := ""
	jsonTag := ""
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		sv := v.FieldByName(sf.Name)
		sfName = sf.Name
		//获取tag
		tag := strings.Split(sf.Tag.Get("json"),",") //获取json tag
		defValue := sf.Tag.Get("default") //获取 default tag
		isRequire := cast.ToBool(sf.Tag.Get("require")) //获取 require tag
		if tag[0] == "" { tag[0] = sf.Name } // 如果json tag为空，则使用字段名
		if tag[0] == "-" { continue }  //如果不转换，则跳过
		jsonTag = tag[0]
		// 从json中获取值
		var value interface{} = nil
		if data != nil {
			value, _ = data[tag[0]]
			// 执行动态动态脚本, 或设置默认值
			if value, err = jsFunc(value); err != nil{
				err = fmt.Errorf("%s(%s) JS脚本执行错误：%s",sf.Name,jsonTag, err.Error())
				return
			}
		}

		if isRequire{
			if value == nil || cast.ToString(value) == ""{
				err = fmt.Errorf("%s(%s) 为必填项", sf.Name, jsonTag)
				return
			}
		}


		if sf.Anonymous { // 匿名字段
			value = data
		}
		if err = reflectSetValue(sf.Type, sv, value, defValue); err != nil {
			err = fmt.Errorf("%s(%s) ：%s", sfName, jsonTag, err.Error())
			return
		}
	}
	defer func(){
		if r := recover(); r != nil{
			err =  fmt.Errorf(fmt.Sprintf("%s(%s) : %s",sfName,jsonTag, r))
		}
	}()

	return
}


func reflectSetValue(rt reflect.Type, rv reflect.Value, value interface{}, defValue string) (err error){
	if value == nil {
		if defValue == ""{
			return
		}else{
			value = defValue
		}
	}
	isPtr := false
	if rt.Kind() == reflect.Ptr {
		isPtr = true
		rt = rt.Elem()
		//if value == nil && defValue != ""{
		//	value = defValue
		//}
	}
	switch rt.Kind() {
	case reflect.Interface: //接口类型
		rv.Set(reflect.ValueOf(value))
	case reflect.String: //字符串型
		tempV := cast.ToString(value)
		if tempV == "" && defValue != "" {
			tempV = defValue
		}
		if isPtr{
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetString(tempV)
		}
	case reflect.Bool: //布尔型
		if cast.ToString(value) == "" && defValue != ""{
			value = defValue
		}
		tempV := cast.ToBool(value)
		if isPtr {
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetBool(tempV)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //整型
		tempV := cast.ToInt64(value)
		if tempV == 0 && defValue != "" {
			tempV = cast.ToInt64(defValue)
		}
		if isPtr{
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetInt(tempV)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //正整型
		tempV := cast.ToUint64(value)
		if tempV == 0 && defValue != "" {
			tempV = cast.ToUint64(defValue)
		}
		if isPtr{
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetUint(tempV)
		}
	case reflect.Float32, reflect.Float64: //浮点型
		tempV := cast.ToFloat64(value)
		if tempV == 0 && defValue != "" {
			tempV = cast.ToFloat64(defValue)
		}
		if isPtr{
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetFloat(tempV)
		}
	case reflect.Struct:
		obj := reflect.New(rt)
		if err = reflectSetStruct(obj.Interface(), value.(map[string]interface{})); err != nil{
			return
		}
		if isPtr{
			rv.Set(obj)
		}else{
			rv.Set(obj.Elem())
		}
	case reflect.Slice, reflect.Array:	//数组、切片
		kind := rt.Elem().Kind()
		arrayValue := value.([]interface{})
		es := make([]reflect.Value,0)
		for _, v := range arrayValue {
			obj := reflect.New(rt.Elem())
			if kind == reflect.Struct || kind == reflect.Map {
				if err = reflectSetStruct(obj.Interface(), v.(map[string]interface{})); err != nil {
					return
				}
			}else{
				if err = reflectSetValue(rt.Elem(), obj.Elem(), v,""); err != nil{
					return
				}
			}
			es = append(es, obj.Elem())
		}
		rv.Set(reflect.Append(rv, es...))
	case reflect.Map: //map
		res := reflect.MakeMap(rt)
		mapV := value.(map[string]interface{})
		for key, v := range mapV{
			k := reflect.ValueOf(key)
			obj := reflect.New(rt.Elem())
			if err = reflectSetValue(rt.Elem(), obj.Elem(), v, ""); err != nil{
				return
			}
			res.SetMapIndex(k, obj.Elem())
		}
		rv.Set(res)
	}
	return
}

// 执行动态动态脚本, 或设置默认值
func jsFunc(value interface{}) (v interface{}, err error){
	if value == nil || reflect.TypeOf(value).Kind() != reflect.String {
		return
	}
	v = value
	vString := value.(string)
	if len(vString) > 3 && strings.ToUpper(vString[:3]) == "JS:" {
		script := fmt.Sprintf("function callfun(search,auth){%s}", string(vString[3:]))
		//script := string(vString[3:])
		vm := goja.New()
		search := vm.NewObject()
		_ = search.Set("a",2)
		auth := vm.NewObject()
		_ = auth.Set("b",3)
		vm.Set("search",search)
		vm.Set("auth",auth)
		if v, err = vm.RunString(script); err != nil{
			return
		}
		if callFun, ok := goja.AssertFunction(vm.Get("callfun")); !ok {
			err = fmt.Errorf("Not a js function")
		}else {
			v, err = callFun(goja.Undefined(), vm.ToValue(search), vm.ToValue(auth))
		}
	}
	return
}