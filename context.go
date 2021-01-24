package mc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strings"
)

type RenderType string

const (
	RenderTypeDefault 		RenderType 	= ""
	RenderTypeHTML 			RenderType	= "HTML"
	RenderTypeJSON 			RenderType 	= "JSON"
	RenderTypeIndentedJSON 	RenderType 	= "INDENTEDJSON"
	RenderTypeSecureJSON 	RenderType	= "SECUREJSON"
	RenderTypeJSONP			RenderType	= "JSONP"
	RenderTypeAsciiJSON		RenderType	= "ASCIIJSON"
	RenderTypeYAML 			RenderType	= "YAML"
	RenderTypeXML			RenderType	= "XML"
	RenderTypeProtoBuf		RenderType	= "PROTOBUF"
	RenderTypeString		RenderType	= "STRING"
)

//结果体，当做为错误panic
type Result struct {
	HttpStatus int                    //响应代码
	Code       string                 //消息代码
	Message    string                 //消息
	Data       interface{}            //数据体
	ExtraData  map[string]interface{} //附加的数据
	RenderType RenderType             //渲染方式
}


func (result *Result) Error() string {
	return result.Message
}


//上下文
type Context struct {
	ModuleName         string     //模块
	RealModuleName     string     //实际使用的模块
	ControllerName     string     //控制器
	RealControllerName string     //实际使用的控制器
	ActionName         string     //操作方法
	ModelName          string     //模型
	RenderType         RenderType //结果渲染类型
	isAjax             bool       //是否ajax提交
	*gin.Context
}

func NewContext(c *gin.Context) (ctx *Context) {
	//去除伪静态后缀
	path := strings.Replace(strings.ToLower(c.Request.URL.Path), "."+option.Router.UrlHtmlSuffix, "", -1)
	//拆分请求路径
	params := strings.Split(path, option.Router.UrlPathSep)
	params = append(params, "", "", "", "")
	ctx = &Context{
		ModuleName:         strings.ToLower(params[1]),
		RealModuleName:     strings.ToLower(params[1]),
		ControllerName:     strings.ToLower(params[2]),
		RealControllerName: strings.ToLower(params[2]),
		ActionName:         strings.ToLower(params[3]),
		ModelName:          strings.ToLower(params[4]),
		isAjax:             c.GetHeader("X-Requested-With") == "XMLHttpRequest",
		Context:            c,
	}

	return
}

// 是否ajax请求
func (ctx *Context) IsAjax() bool {
	return ctx.isAjax
}

//获取分页信息
func (ctx *Context) getPage() (page int, pageSize int) {
	if ctx.Request.Method == "GET" {
		page = cast.ToInt(ctx.DefaultQuery(option.Request.PageName, "1"))
		pageSize = cast.ToInt(ctx.DefaultQuery(option.Request.PageSizeName, "50"))
	} else {
		page = cast.ToInt(ctx.DefaultPostForm(option.Request.PageName, "1"))
		pageSize = cast.ToInt(ctx.DefaultPostForm(option.Request.PageSizeName, "50"))
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = option.Request.PageSizeValue
	}
	return
}

//获取排序信息
func (ctx *Context) getOrder() (order string) {
	if ctx.Request.Method == "GET" {
		order = ctx.DefaultQuery(option.Request.OrderName, "")
	} else {
		order = ctx.DefaultPostForm(option.Request.OrderName, "")
	}
	return
}

func (ctx *Context) Render(renderType RenderType, httpStatus int, template string, assign *Assign) {
	//渲染类型判断 (如果未指定，则使用上下文，如果上下文没指定，则判断是否为ajax提交)
	if renderType == RenderTypeDefault {
		if ctx.RenderType == RenderTypeDefault && ctx.isAjax {
			renderType = option.Response.AjaxRenderType
		} else {
			renderType = ctx.RenderType
		}
	}

	switch renderType {
	case RenderTypeAsciiJSON:
		ctx.AsciiJSON(httpStatus, assign.Result)
	case RenderTypeIndentedJSON:
		ctx.IndentedJSON(httpStatus, assign.Result)
	case RenderTypeJSON:
		ctx.JSON(httpStatus, assign.Result)
	case RenderTypeJSONP:
		ctx.JSONP(httpStatus, assign.Result)
	case RenderTypeSecureJSON:
		ctx.SecureJSON(httpStatus, assign.Result)
	case RenderTypeYAML:
		ctx.YAML(httpStatus, assign.Result)
	case RenderTypeXML:
		ctx.XML(httpStatus, assign.Result)
	case RenderTypeProtoBuf:
		ctx.ProtoBuf(httpStatus, assign.Result)
	case RenderTypeString:
		ctx.String(httpStatus, fmt.Sprint(assign.Result))
	default:
		ctx.HTML(httpStatus, template, assign)
	}
}