package mc

import (
    "bytes"
    "fmt"
    "github.com/gin-gonic/gin/render"
    "github.com/spf13/cast"
    "html/template"
    "strings"
)



var TemplateFuncMap = template.FuncMap{
    "html":HtmlUnescaped,
    "string": HtmlString,
    "int": HtmlInt,
    "float": HtmlFloat,
    "bool": htmlBool,
}

type widgetAssign struct{
    Field *ModelBaseField
    Value interface{}
    FromKvs Kvs
}


func HtmlUnescaped(x string) interface{}{
    return template.HTML(x)
}

func HtmlString(x interface{}) string{
    return cast.ToString(x)
}

func HtmlInt(x interface{}) int{
    return cast.ToInt(x)
}

func HtmlFloat(x interface{}) float32{
    return cast.ToFloat32(x)
}

func htmlBool(x interface{}) bool {
    return cast.ToBool(x)
}


//widget 的html
func CreateWidget(field *ModelBaseField, value interface{}, fromKvs Kvs) string {
    var buf bytes.Buffer
    key := strings.ToLower(field.Widget)
    if key == "" {
        key = "text"
    }
    key = fmt.Sprintf("widget/%s", key)
    widgetAssign := widgetAssign{Field: field, Value: value, FromKvs: fromKvs}
    globalTemplate := option.engine.HTMLRender.Instance("", nil).(render.HTML)
    if err := globalTemplate.Template.ExecuteTemplate(&buf, key, widgetAssign); err != nil {
        panic(fmt.Errorf("%s(%s) widget 渲染失败：%s", field.Name, key, err))
    }
    return buf.String()
}
