package mc


//type (
//	ErrorCallBackFunc func(httpStatus int, ctx *Context, err error) map[string]interface{}
//	SuccessCallBackFunc func(ctx *Context, result *Result) interface{}
//)

// 默认错误结果回调
//func defaultErrorCallBackFunc (httpStatus int, ctx *Context, err error) map[string]interface{} {
//	return map[string]interface{}{
//		"code": httpStatus,
//		"message": err.Error(),
//	}
//}

// 默认成功结果回调
//func defaultSuccessCallBackFunc (ctx *Context, result *Result) interface{} {
//	return result
//}


//mc选项结构体
type Option struct {
	DefaultConnName            string                            //默认数据库连接名			默认值："default"
	PageSize                   int                               //默认页记录数				默认值：50
	VarPage                    string                            //前端页码参数名				默认值："page"
	VarPageSize                string                            //前端页记录数参数名			默认值："page_size"
	UrlPathSep                 string                            //URL路径之间的分割符号（不能使用_下线线）	默认为 "/"
	UrlHtmlSuffix              string                            //URL伪静态后缀设置			默认为 "html"
	Controllers                map[string]map[string]IController //控制器map
	NotAuthActions             []string                          //无须登录认证的操作方法列表
	NotAuthRedirect            string                            //未登录跳转到页面地址 		默认值 "/admin/public/login.html"
	BaseModuleMapKey           string                            //全局基础模块key			默认值 "base"
	BaseControllerMapKey       string                            //全局基础控制器key			默认值 "base"
	ModuleBaseControllerMapKey string                            //当前模块下基础控制器key	默认为 "base"
	//ErrorCallBackFunc          ErrorCallBackFunc                 //错误回调
	//SuccessCallBackFunc        SuccessCallBackFunc               //成功回调
	ErrorTemplate              string                            //错误页面模版				默认值 "error.html"
	ModelConfigsFilePath       string                            //自定义模型配置文件存放路径	默认值："./model_configs/"
	widgetTemplatePath         string                 			//小物件模版					默认值："./widgets/"
}

//默认选项设置
var option = Option{
	DefaultConnName:            "default",
	VarPage:                    "page",
	VarPageSize:                "page_size",
	PageSize:                   50,
	UrlPathSep:                 "/",
	UrlHtmlSuffix:              "html",
	Controllers:                map[string]map[string]IController{},
	NotAuthActions:             []string{},
	NotAuthRedirect:            "/public/login.html",
	BaseModuleMapKey:           "base",
	BaseControllerMapKey:       "base",
	ModuleBaseControllerMapKey: "base",
	//ErrorCallBackFunc:          defaultErrorCallBackFunc,
	//SuccessCallBackFunc:        defaultSuccessCallBackFunc,
	ErrorTemplate:              "error.html",
	ModelConfigsFilePath:       "./model_configs/",
	widgetTemplatePath:			"./widgets/",
}

func Default() *Option {
	return &option
}

func (o *Option) WidgetTemplatePath (path string){
	o.widgetTemplatePath = path
	initWidgets()
}