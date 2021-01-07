package mc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strings"
)

const (
	RenderTypeHTML 			= "HTML"
	RenderTypeJSON 			= "JSON"
	RenderTypeIndentedJSON 	= "INDENTEDJSON"
	RenderTypeSecureJSON 	= "SECUREJSON"
	RenderTypeJSONP			= "JSONP"
	RenderTypeAsciiJSON		= "ASCIIJSON"
	RenderTypeYAML 			= "YAML"
	RenderTypeXML			= "XML"
	RenderTypeProtoBuf		= "PROTOBUF"
	RenderTypeString		= "STRING"
)



//上下文
type Context struct {
	ModuleName         string      //模块
	RealModuleName     string      //实际使用的模块
	ControllerName     string      //控制器
	RealControllerName string      //实际使用的控制器
	ActionName         string      //操作方法
	ModelName          string      //模型
	Page               int         //页码
	PageSize           int         //记录数
	RenderType         string      //结果渲染类型
	Result             interface{} //数据内容
	*gin.Context
}

func NewContext(c *gin.Context) (ctx *Context) {
	//去除伪静态后缀
	path := strings.Replace(strings.ToLower(c.Request.URL.Path), "."+option.UrlHtmlSuffix, "", -1)
	//拆分请求路径
	params := strings.Split(path, option.UrlPathSep)
	params = append(params, "", "", "", "")
	ctx = &Context{
		ModuleName:         strings.ToLower(params[1]),
		RealModuleName:     strings.ToLower(params[1]),
		ControllerName:     strings.ToLower(params[2]),
		RealControllerName: strings.ToLower(params[2]),
		ActionName:         strings.ToLower(params[3]),
		ModelName:          strings.ToLower(params[4]),
		Page:               0,
		PageSize:           0,
		Context:            c,
	}
	//页码处理
	ctx.Page, ctx.PageSize = ctx.getPage()
	return
}


//获取分页信息
func (ctx *Context) getPage() (page int, pageSize int) {
	if ctx.Request.Method == "GET" {
		page = cast.ToInt(ctx.Query(option.VarPage))
		pageSize = cast.ToInt(ctx.Query(option.VarPageSize))
	} else {
		page = cast.ToInt(ctx.Request.Form[option.VarPage])
		page = cast.ToInt(ctx.Request.Form[option.VarPageSize])
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = option.PageSize
	}
	return
}

func (ctx *Context) Render(httpStatus int, template string, result interface{}) {
	switch strings.ToUpper(ctx.RenderType) {
	case RenderTypeAsciiJSON:
		ctx.AsciiJSON(httpStatus, result)
	case RenderTypeIndentedJSON:
		ctx.IndentedJSON(httpStatus, result)
	case RenderTypeJSON:
		ctx.JSON(httpStatus, result)
	case RenderTypeJSONP:
		ctx.JSONP(httpStatus, result)
	case RenderTypeSecureJSON:
		ctx.SecureJSON(httpStatus, result)
	case RenderTypeYAML:
		ctx.YAML(httpStatus, result)
	case RenderTypeXML:
		ctx.XML(httpStatus, result)
	case RenderTypeProtoBuf:
		ctx.ProtoBuf(httpStatus, result)
	case RenderTypeString:
		ctx.String(httpStatus, fmt.Sprint(result))
	default:
		ctx.HTML(httpStatus, template, result)
	}
}