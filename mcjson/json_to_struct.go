package mcjson

import (
	"github.com/spf13/cast"
	"reflect"
	"strings"
)


func JsonToStruct(jsonData []byte, obj interface{}) error{
	data := Json(jsonData).Getdata()
	return reflectSetStruct(obj, data)
}

func toStruct(data map[string]interface{}, out interface{}) (err error){
	//defer func(){
	//	if e := recover(); e != nil {
	//		err = fmt.Errorf("%s",e)
	//	}
	//}()

	v := reflect.ValueOf(out).Elem()
	t := reflect.TypeOf(out).Elem()



	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		sv := v.FieldByName(sf.Name)

		jsonTag := strings.Split(sf.Tag.Get("json"),",") //获取json tag
		def := sf.Tag.Get("default") //获取 default tag
		if jsonTag[0] == "" { jsonTag[0] = sf.Name } // 如果json tag为空，则使用字段名
		if jsonTag[0] == "-" { continue }  //如果不转换，则跳过
		// 从json中获取值
		var value interface{} = nil
		if data != nil {
			ok := false
			if value, ok = data[jsonTag[0]]; !ok {
				value = nil
			}
		}

		// 判断值是否是动态脚本
		if value != nil {
			if reflect.TypeOf(value).Kind() == reflect.String {
				vString := value.(string)
				if len(vString) > 6 && vString[:6] == "eveal:" {
					value = vString
				}
			}
		} else if def != "" {
			value = def
		}


		switch sf.Type.Kind() {
		case reflect.String: //字符串型
			sv.SetString(cast.ToString(value))

		case reflect.Bool: //布尔型
			sv.SetBool(cast.ToBool(value))

		case reflect.Int,reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //整型
			sv.SetInt(cast.ToInt64(value))

		case reflect.Float32, reflect.Float64: //浮点型
			sv.SetFloat(cast.ToFloat64(value))

		case reflect.Struct: //结构
			newStruct := reflect.New(sf.Type)
			if sf.Anonymous { // 匿名字段
				err = toStruct(data, newStruct.Interface())
			}else if value == nil {
				err = toStruct(nil, newStruct.Interface())
			}else{
				err = toStruct(value.(map[string]interface{}), newStruct.Interface())
			}
			if err != nil {
				return
			}
			sv.Set(newStruct.Elem())

		case reflect.Ptr: //指针
			tx := reflect.TypeOf(sv.Interface())
			switch tx.Elem().Kind() {
			case reflect.String:
				tempV := cast.ToString(value)
				sv.Set(reflect.ValueOf(&tempV))
			case reflect.Bool:
				tempV := cast.ToBool(value)
				sv.Set(reflect.ValueOf(&tempV))
			case reflect.Int,reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				tempV := cast.ToInt64(value)
				sv.Set(reflect.ValueOf(&tempV))
			case reflect.Float32, reflect.Float64:
				tempV := cast.ToFloat64(value)
				sv.Set(reflect.ValueOf(&tempV))
			case reflect.Struct:
				if sv.IsNil() {
					newStruct := reflect.New(tx.Elem())
					sv.Set(newStruct)
				}
				if sf.Anonymous { // 匿名字段
					err = toStruct(data, sv.Interface())
				}else if value == nil {
					err = toStruct(nil,sv.Interface())
				}else{
					err = toStruct(value.(map[string]interface{}), sv.Interface())
				}
				if err != nil {
					return
				}
			}
		case reflect.Slice, reflect.Array:	//数组、切片
			if value == nil { return }
			kind := sf.Type.Elem().Kind()
			elem := sf.Type.Elem()
			arrayValue := value.([]interface{})

			es := make([]reflect.Value,0)

			for _, v := range arrayValue {
				if kind == reflect.Struct || kind == reflect.Map {
					newStruct := reflect.New(elem)
					if err = toStruct(v.(map[string]interface{}), newStruct.Interface()); err != nil {
						return
					}
					es = append(es, newStruct.Elem())
				}else{

					es = append(es, reflect.ValueOf(v))
				}
			}
			sv.Set(reflect.Append(sv, es...))
		case reflect.Map:  //map
			if value == nil { return }
			kind := sf.Type.Elem().Kind()
			elem := sf.Type.Elem()
			res := reflect.MakeMap(sf.Type)
			mapV := value.(map[string]interface{})
			for key, v := range mapV{
				k := reflect.ValueOf(key)

				switch kind {
				case reflect.Interface:
					res.SetMapIndex(k,reflect.ValueOf(v))
				case reflect.Struct, reflect.Map:
					newObj := reflect.New(elem)
					if err = toStruct(v.(map[string]interface{}), newObj.Interface()); err != nil {
						return
					}
					res.SetMapIndex(k, newObj.Elem())
				case reflect.Slice, reflect.Array:	//数组、切片

				default:
				}


			}
			sv.Set(res)
		}
	}
	return nil
}

func reflectSetStruct(obj interface{},data map[string]interface{}) (err error){
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		sv := v.FieldByName(sf.Name)

		//获取tag
		tag := strings.Split(sf.Tag.Get("json"),",") //获取json tag
		defValue := sf.Tag.Get("default") //获取 default tag
		if tag[0] == "" { tag[0] = sf.Name } // 如果json tag为空，则使用字段名
		if tag[0] == "-" { continue }  //如果不转换，则跳过

		// 从json中获取值
		var value interface{} = nil
		if data != nil {
			value, _ = data[tag[0]]
		}
		value = jsFunc(value, defValue) // 执行动态动态脚本, 或设置默认值
		if sf.Anonymous { // 匿名字段
			value = data
		}
		switch sf.Type.Kind() {
		//case reflect.Ptr: // 指针
		//	reflectSetPar(sv, value)
		default:
			if err = reflectSetValue(sf.Type, sv, value); err != nil {
				return
			}
		}

	}
	return
}


func reflectSetValue(rt reflect.Type, rv reflect.Value, value interface{}) (err error){
	if value == nil { return }
	isPtr := false
	if rt.Kind() == reflect.Ptr {
		isPtr = true
		rt = rt.Elem()
	}
	switch rt.Kind() {
	case reflect.Interface: //接口类型
		rv.Set(reflect.ValueOf(value))
	case reflect.String: //字符串型
		tempV := cast.ToString(value)
		if isPtr{
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetString(tempV)
		}
	case reflect.Bool: //布尔型
		tempV := cast.ToBool(value)
		if isPtr {
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetBool(tempV)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //整型
		tempV := cast.ToInt64(value)
		if isPtr{
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetInt(tempV)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //正整型
		tempV := cast.ToUint64(value)
		if isPtr{
			rv.Set(reflect.ValueOf(&tempV))
		}else {
			rv.SetUint(tempV)
		}
	case reflect.Float32, reflect.Float64: //浮点型
		tempV := cast.ToFloat64(value)
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
				if err = reflectSetValue(rt.Elem(), obj.Elem(), v); err != nil{
					return
				}
			}
			es = append(es, obj.Elem())
		}
		rv.Set(reflect.Append(rv, es...))
	case reflect.Map: //map
		//kind := sf.Type.Elem().Kind()
		//elem := sf.Type.Elem()
		res := reflect.MakeMap(rt)
		mapV := value.(map[string]interface{})
		for key, v := range mapV{
			k := reflect.ValueOf(key)
			obj := reflect.New(rt.Elem())
			if err = reflectSetValue(rt.Elem(), obj.Elem(), v); err != nil{
				return
			}
			res.SetMapIndex(k, obj.Elem())
		}
		rv.Set(res)
	}
	return
}

// 执行动态动态脚本, 或设置默认值
func jsFunc(value interface{}, defValue interface{}) interface{}{
	if value != nil {
		if reflect.TypeOf(value).Kind() == reflect.String {
			vString := value.(string)
			if len(vString) > 3 && vString[:3] == "JS:" {
				value = vString
			}
		}
	} else if defValue != "" {
		value = defValue
	}
	return value
}