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
    "Html":HtmlUnescaped,
    "String": HtmlString,
    "Int": HtmlInt,
    "Float": HtmlFloat,
    "Bool": htmlBool,
    "Split": strings.Split,
    "Join": strings.Join,
    "InArray": InArray,
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

//判断棋模版是否存在
func HasTemplate(name string) bool{
    globalTemplate :=option.engine.HTMLRender.Instance("", nil).(render.HTML)
    for _, v := range globalTemplate.Template.Templates(){
        if name == v.Name(){
            return true
        }
    }
    return false
}

func GetTemplateName(ctx *Context) (name string){
    //先查找路由指定的模板
    name = fmt.Sprintf("%s/%s_%s.html", ctx.ModuleName, ctx.ControllerName, ctx.ActionName)
    if HasTemplate(name) {
        return
    }
    //找不到，再查找路由指定的模块下公共控制器对应模板
    name = fmt.Sprintf("%s/%s_%s.html", ctx.ModuleName, option.Router.ModuleBaseControllerName, ctx.ActionName)
    if HasTemplate(name){
        return
    }

    //再找不到，使用全局基础模板
    name = fmt.Sprintf("%s/%s_%s.html", option.Router.BaseModuleName, option.Router.BaseControllerName, ctx.ActionName)
    return
}
