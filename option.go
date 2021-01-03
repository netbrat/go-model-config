package mc


import (
	"github.com/gin-gonic/gin"
)

type RouterMap struct {
	Controllers    map[string]map[string]IController //控制器map
	NotAuthActions []string                          //无须登录认证的操作方法列表
}

//mc选项结构体
type Option struct {
	PageSize                   int                                            //默认页记录数				默认值：50
	ErrorTemplate              string                                         //错误页面模版				默认值 "error.html"
	RouterMap                  RouterMap                                      //路由表
	NotAuthRedirect            string                                         //未登录跳转到页面地址 		默认值 "/admin/public/login.html"
	UrlPathSep                 string                                         //URL路径之间的分割符号（不能使用_下线线）	默认为 "/"
	UrlHtmlSuffix              string                                         //URL伪静态后缀设置			默认为 "html"
	BaseModuleMapKey           string                                         //全局基础模块key			默认值 "base"
	BaseControllerMapKey       string                                         //全局基础控制器key			默认值 "base"
	ModuleBaseControllerMapKey string                                         //当前模块下基础控制器key	默认为 "base"
	ErrorCallBackFunc          func(httpStatus int, error string) interface{} //错误回调
	VarPage                    string                                         //前端页码参数名			默认值："page"
	VarPageSize                string                                         //前端页记录数参数名		默认值："page_size"
	ConfigsFilePath            string                                         //自定义模型配置文件存放路径	默认值："./model_configs/"
	DefaultConnName            string                                         //默认数据库连接名			默认值："default"
}

//默认选项设置
var option = Option{
	PageSize:                   50,
	ErrorTemplate:              "error.html",
	RouterMap:                  RouterMap{},
	NotAuthRedirect:            "/public/login.html",
	UrlPathSep:                 "/",
	UrlHtmlSuffix:              "html",
	BaseModuleMapKey:           "base",
	BaseControllerMapKey:       "base",
	ModuleBaseControllerMapKey: "base",
	ErrorCallBackFunc:          func(httpStatus int, error string) interface{} { return gin.H{"code": httpStatus, "message": error} },
	VarPage:                    "page",
	VarPageSize:                "page_size",
	ConfigsFilePath:            "./model_configs/",
	DefaultConnName:            "default",
}

func Default() *Option {
	return &option
}

