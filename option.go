package mc

import "github.com/gin-gonic/gin"

//路由列表结构体
type RouterMap struct {
	Controllers    map[string]map[string]IController
	NotAuthActions []string
}


//选项结构体
type Option struct {
	PageSize        			int       	//默认页记录数				默认值：50
	ErrorTemplate   			string    	//错误页面模版				默认值 "error.html"
	RouterMap       			RouterMap 	//路由列表					默认值 nil
	NotAuthRedirect 			string    	//未登录跳转到页面地址 		默认值 "/admin/public/login.html"
	UrlPathSeparator			string		//URL路径之间的分割符号（不能使用_下线线）	默认为 "/"
	UrlHtmlSuffix 				string		//URL伪静态后缀设置			默认为 "html"
	BaseModuleMapKey			string		//全局基础模块key			默认值 "base"
	BaseControllerMapKey		string		//全局基础控制器key			默认值 "base"
	ModuleBaseControllerMapKey	string		//当前模块下基础控制器key		默认为 "base"
	ErrorCallBackFunc			func(httpStatus int, error string) interface{}	//错误回调
	VarPage						string		//前端页码参数名				默认值："page"
	VarPageSize					string		//前端页记录数参数名			默认为："page_size"
}


//默认选项设置
var option = Option{
	PageSize:                   50,
	ErrorTemplate:              "error.html",
	RouterMap:                  RouterMap{},
	NotAuthRedirect:            "/admin/public/login.html",
	UrlPathSeparator:           "/",
	UrlHtmlSuffix:              "html",
	BaseModuleMapKey:           "base",
	BaseControllerMapKey:       "base",
	ModuleBaseControllerMapKey: "base",
	ErrorCallBackFunc: 			func(httpStatus int, error string) interface{}{ return gin.H{"code":httpStatus, "message":error}},
	VarPage:					"page",
	VarPageSize:				"page_size",
}

func DefaultOption() *Option{
	return &option;
}
