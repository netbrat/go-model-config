package mc

import (
    "bytes"
    "fmt"
    "github.com/gin-gonic/gin/render"
    "github.com/spf13/cast"
    "html/template"
    "strings"
)


// 默认的模版
var defaultWidgetsTemplate map[string]string = map[string]string{
    "widget/text": `<input type="text" id="{{.Field.Name}}" name="{{.Field.Name}}" value="{{.Value}}"/>`,
    "widget/textarea": `<textarea id="{{.Field.Name}}" name="{{.Field.Name}}">{{.Value}}</textarea>`,
    "widget/select": `<select id="{{.Field.Name}}" name="{{.Field.Name}}"><option value=""></option></select>`,
    "widget/date": `<input type="date" id="{{.Field.Name}}" name="{{.Field.Name}}" value="{{.Value}}" />`,
}

var TemplateFuncMap = template.FuncMap{
    "HtmlUnescaped":HtmlUnescaped,
    "string": HtmlString,
}

var globalTemplate render.HTML


type widgetAssign struct{
    Field *ModelBaseField
    Value interface{}
    FromData Enum
}


func HtmlUnescaped(x string) interface{}{
    return template.HTML(x)
}

func HtmlString(x interface{}) string{
    return cast.ToString(x)
}


// 初始化所有 widget 模版
func initWidgets() {
    globalTemplate = option.engine.HTMLRender.Instance("", nil).(render.HTML)
    tmplNames := make([]string, 0)
    for _, templ := range globalTemplate.Template.Templates() {
        tmplNames = append(tmplNames,templ.Name())
    }
    for key, content := range defaultWidgetsTemplate {
        if !inArray(key,tmplNames){
            _, err := globalTemplate.Template.New(key).Funcs(option.engine.FuncMap).Parse(content)
            if err != nil {
                panic(fmt.Errorf("默认模版[%]解析失败：%s", key,err))
            }
            //option.engine.SetHTMLTemplate(templ)
        }
    }
}

//widget 的html
func CreateWidget(field *ModelBaseField, value interface{}, FromData Enum) string {
    var buf bytes.Buffer
    key := strings.ToLower(field.Widget)
    if key == "" { key = "text" }
    key = fmt.Sprintf("widget/%s", key)
    widgetAssign := widgetAssign{Field: field, Value: value, FromData:  FromData}
    //html := option.engine.HTMLRender.Instance("", widgetAssign).(render.HTML)

    if err :=  globalTemplate.Template.ExecuteTemplate(&buf,key,widgetAssign); err != nil {
        panic(fmt.Errorf("%s(%s) widget 渲染失败：%s",field.Name, key,err))
    }
    return buf.String()
}
