package mc

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strings"
	"unicode"
)

//驼峰转下划线小写
func CamelToCase(str string) string{
	s := []rune(str)
	desc := ""
	for i, v := range s {
		if unicode.IsUpper(v) {
			if i > 0 {
				desc += "_"
			}
		}
		desc += fmt.Sprintf("%c", v)
	}
	return strings.ToLower(desc)
}

// 转大驼峰写法
func CaseToUpperCamel(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

//转小驼峰
func CaseToLowerCamel(str string) string{
	return Lcfirst(CaseToUpperCamel(str))
}


// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}


// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}


//查找元素是否在数组中
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}


func getOffsetLimit(page int, pageSize int)(offset int, limit int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = option.Request.PageSizeValue
	}
	limit = pageSize
	offset = (page - 1) * pageSize
	return
}

//生成一串n个相同字符串组成的的字符串
func nString(str string, n int)(s string) {
	s = ""
	for i := 0; i < n; i++ {
		s += str
	}
	return
}


func ToTreeMap(source []map[string]interface{}, keyName string) (desc map[string]interface{}) {
	desc = toTreeMap(source, keyName, 0, "")
	return
}


func toTreeMap(source []map[string]interface{}, keyName string, startIndex int, parent string) (desc map[string]interface{}){
	desc = make(map[string]interface{})
	for i := startIndex; i < len(source); i++ {
		nodeChildCount := cast.ToInt(source[i]["__mc_child_count"])
		nodeParent := cast.ToString(source[i]["__mc_parent"])
		nodeKey := cast.ToString(source[i][keyName])
		if nodeParent == parent {
			newNode := source[i]
			if nodeChildCount > 0 {
				newNode["__mc_children"] = toTreeMap(source, keyName, i+1, nodeKey)
			}else{
				newNode["__mc_children"] = nil
			}
			desc[nodeKey] = newNode
		}
	}
	return
}