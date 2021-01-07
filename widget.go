package mc

import (
    "bytes"
    "fmt"
    "html/template"
    "io/ioutil"
    "path"
    "strings"
)

// 默认的模版
var defaultWidgetsTemplate map[string]string = map[string]string{
    "text": `<input type="text" id="{{.Field.Name}}" name="{{.Field.Name}}" value="{{.Value}}"/>`,
    "textarea": `<textarea id="{{.Field.Name}}" name="{{.Field.Name}}">{{.Value}}</textarea>`,
    "select": `<select id="{{.Field.Name}}" name="{{.Field.Name}}"><option value=""></option></select>`,
    "date": `<input type="date" id="{{.Field.Name}}" name="{{.Field.Name}}" value="{{.Value}}" />`,
}

var widgets = make(map[string] *template.Template)

type WidgetAssign struct{
    Field *ModelBaseField
    Value interface{}
    FromData map[string]interface{}
}

// 初始化所有 widget 模版
func initWidgets() {
    var err error
    for key, tpl := range defaultWidgetsTemplate {
        widgets[key], err = template.New(key).Parse(tpl)
        if err != nil{
            panic(fmt.Errorf("widget 默认模版[%]解析失败：%s", key,err))
        }
    }

    if option.widgetTemplatePath == "" { return }

    if rd, err := ioutil.ReadDir(option.widgetTemplatePath); err != nil {
        panic(fmt.Errorf("widget模版目录[%s]读取失败：%s", option.widgetTemplatePath, err))
    }else {
        for _, file := range rd {
            fileName := file.Name()
            if !file.IsDir() && strings.ToLower(path.Ext(fileName)) ==".html" {
                fileBase := strings.ToLower(fileName[:len(fileName)-5])
                if widgets[fileBase], err = template.New(fileBase).ParseFiles(option.widgetTemplatePath + fileName); err != nil {
                    panic(fmt.Errorf("widget模版文件[%s%s]解析失败：%s", option.widgetTemplatePath, fileName, err))
                }else {
                }
            }
        }
    }
}

//widget 的html
func CreateWidget(field *ModelBaseField, value interface{}, fromData map[string]interface{}) string {
    var buf bytes.Buffer
    key := strings.ToLower(field.Widget)
    if key == "" {
        key = "text"
    }
    if _, ok := widgets[key]; !ok{
        return ""
    }
    widgetAssign := WidgetAssign{Field: field, Value: value, FromData:  fromData}
    if err := widgets[key].ExecuteTemplate(&buf,key,widgetAssign); err != nil {
        panic(err)
    }

    return buf.String()
}
