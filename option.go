package mc

import "github.com/gin-gonic/gin"


// 调函数定义
type (
	ContextBeforeFunc func(c *gin.Context) (err error) // 初始化上下文前回调
	ContextAfterFunc func(ctx *Context) (err error) //初始化上下文后回调
	GetRowAuthModelsFunc func() []string //获取所有具有行权限模型回调
	GetRowAuthFunc func(modelName string) (authCodes []string, isAll bool) //获取行权限回调
	GetColAuthFunc func(modelName string) (authCodes []string, isAll bool) //获取列权限回调
)

//权限选项
type authOption struct {
	GetRowAuthModelsCallback GetRowAuthModelsFunc
	GetRowAuthCallback GetRowAuthFunc
	GetColAuthCallback GetColAuthFunc
}

//响应选项
type responseOption struct {
	CodeName         string     //代码项的key
	MessageName      string     //消息项的key
	DataName         string     //数据项的key
	TotalName        string     //总记录数或影响的记录数项的key
	FootName         string     //表尾汇总数据项的key
	SuccessCodeValue string     //成功代码值
	FailCodeValue    string     //失败默认代码值
	AjaxRenderType   RenderType //默认ajax渲染类型
}

//请求选项
type requestOption struct {
	OrderName             string            //排序字段  				默认值 order
	PageName              string            //前端页码参数名				默认值："page"
	PageSizeName          string            //前端页记录数参数名			默认值："limit"
	PageSizeValue         int               //默认页记录数				默认值：50
	ContextBeforeCallback ContextBeforeFunc //初始化上下文前回调
	ContextAfterCallback  ContextAfterFunc  //初始化上下文后回调
}

//路由选项
type routerOption struct {
	UrlPathSep               string                            //URL路径之间的分割符号（不能使用_下线线）	默认为 "/"
	UrlHtmlSuffix            string                            //URL伪静态后缀设置	默认为 "html"
	ControllerMap            map[string]map[string]IController //控制器map
	BaseModuleName           string                            //全局基础模块key	默认值 "base"
	BaseControllerName       string                            //全局基础控制器key	默认值 "base"
	ModuleBaseControllerName string                            //当前模块下基础控制器key	默认为 "base"
}


//mc选项结构体
type Option struct {
	engine               *gin.Engine    //
	DefaultConnName      string         //默认数据库连接名			默认值："default"
	ErrorTemplate        string         //错误页面模版				默认值 "error.html"
	ModelConfigsFilePath string         //自定义模型配置文件存放路径	默认值："./mconfigs/"
	Router               routerOption   //路由选项
	Response             responseOption //结果项设置
	Request              requestOption  //请求项设置
	ModelAuth            authOption     //权限项设置
}

//默认选项设置
var option = Option{
	DefaultConnName:      "default",
	ErrorTemplate:        "error.html",
	ModelConfigsFilePath: "./mconfigs/",
	ModelAuth: authOption{
		GetRowAuthModelsCallback: func() []string { return make([]string, 0) },
		GetRowAuthCallback:       func(modelName string) (authCodes []string, isAll bool) { return make([]string, 0), true },
		GetColAuthCallback:       func(modelName string) (authCodes []string, isAll bool) { return make([]string, 0), true },
	},
	Router: routerOption{
		UrlPathSep:               "/",
		UrlHtmlSuffix:            "html",
		ControllerMap:            map[string]map[string]IController{},
		BaseModuleName:           "base",
		BaseControllerName:       "base",
		ModuleBaseControllerName: "base",
	},
	Response: responseOption{
		SuccessCodeValue: "0000",
		FailCodeValue:    "1000",
		CodeName:         "code",
		MessageName:      "msg",
		DataName:         "data",
		TotalName:        "total",
		FootName:         "foot",
		AjaxRenderType:   RenderTypeJSON,
	},
	Request: requestOption{
		PageName:              "page",
		PageSizeName:          "limit",
		PageSizeValue:         50,
		OrderName:             "order",
		ContextBeforeCallback: func(c *gin.Context) (err error) { return },
		ContextAfterCallback:  func(ctx *Context) (err error) { return },
	},
}

func Default(engine *gin.Engine) *Option {
	option.engine = engine
	return &option
}


