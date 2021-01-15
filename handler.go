package mc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

//未初始化控制器之前的错误处理
func AbortWithErrorService(c *gin.Context, result Result){
	var obj IController
	controller, exist := c.Get("IController")
	if exist && controller != nil {
		obj = controller.(IController)
	}else{
		obj = &Controller{}
		//初始化上下文
		obj.Initialize(NewContext(c))
	}
	obj.AbortWithError(result)
}


//所有入口适配器
func HandlerAdapt(c *gin.Context) {
	//对post数据进行处理
	if c.ContentType() == "multipart/form-data" {
		_ = c.Request.ParseMultipartForm(1048576)
	} else {
		_ = c.Request.ParseForm()
	}


	//初始化上下文
	ctx := NewContext(c)
	//判断是否需要登录验证
	actionString := fmt.Sprintf("%s.%s.%s", ctx.ModelName, ctx.ControllerName, ctx.ActionName)
	if !inArray(actionString, option.Router.NotAuthActions) {
		//验证是否登录
	}

	//从注册表中查询路由指定的控制器
	var obj IController
	ok := false
	//先查找路由指定的模块及控制器
	if obj, ok = option.Router.ControllerMap[ctx.ModuleName][ctx.ControllerName]; ok {
		//找到不作处理
	} else if obj, ok = option.Router.ControllerMap[ctx.ModuleName][option.Router.ModuleBaseControllerName]; ok {
		//找不到，再查找路由指定的模块下公共控制器
		ctx.RealControllerName = option.Router.ModuleBaseControllerName
	} else if  obj, ok = option.Router.ControllerMap[option.Router.BaseModuleName][option.Router.BaseControllerName]; ok {
		//再查找app下的公共控制器
		ctx.RealModuleName = option.Router.BaseModuleName
		ctx.RealControllerName = option.Router.BaseControllerName
	} else {
		//都找不到，报错误
		msg := fmt.Sprintf("未找到对应的的控制器页面[%s.%s]", ctx.ModuleName, ctx.ControllerName)
		AbortWithErrorService(c, Result{HttpStatus:404, Message:msg})
		return
	}

	//主要是为了发生异常，显示错误时再次初始化一个默认控制器
	c.Set("IController", obj)
	//初始化控制器
	obj.Initialize(ctx)

	//判断控制器内的操作方法是否存在
	//先判断 XxxGet,XxxPost方式，再判断Xxx
	objValue := reflect.ValueOf(obj)
	actionName := fmt.Sprintf("%s%s", ToCamelCase(ctx.ActionName,false), ctx.Request.Method)
	fn := objValue.MethodByName(actionName)
	if fn.Kind() != reflect.Func {
		actionName = ToCamelCase(ctx.ActionName, false)
		fn = objValue.MethodByName(actionName)
		if fn.Kind() != reflect.Func {
			msg := fmt.Sprintf("未找到对应的操作方法[%s.%s.%s]", ctx.ModuleName, ctx.ControllerName, ctx.ActionName)
			obj.AbortWithError(Result{HttpStatus:404, Code:"404", Message:msg})
			return
		}
	}

	//调用操作方法
	fn.Call(nil)
}
