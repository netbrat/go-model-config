package mc

import "html/template"

func HtmlUnescaped(x string) interface{}{
    return template.HTML(x)
}

var TemplateFuncMap = template.FuncMap{
    "HtmlUnescaped":HtmlUnescaped,
}

