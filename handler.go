package mc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

//未初始化控制器之前的错误处理
func AbortWithError(c *gin.Context, httpStatus int, result interface{}){
	//初始化上下文
	ctx := &Context{Context: c}
	ctx.Init()
	obj := &Controller{}
	obj.Init(ctx)
	obj.AbortWithError(http.StatusNotFound, result)
}


//所有入口适配器
func HandlerAdapt(c *gin.Context) {
	//初始化上下文
	ctx := &Context{Context: c}
	ctx.Init()
	//判断是否需要登录验证
	actionString := fmt.Sprintf("%s.%s.%s", ctx.Reqs.Model, ctx.Reqs.Controller, ctx.Reqs.Action)
	if !InArray(actionString, option.RouterMap.NotAuthActions) {
		//验证是否登录
	}

	//从注册表中查询路由指定的控制器
	var obj IController
	ok := false
	//先查找路由指定的模块及控制器
	if obj, ok = option.RouterMap.Controllers[ctx.Reqs.Module][ctx.Reqs.Controller]; ok {
		//找到不作处理
	} else if obj, ok = option.RouterMap.Controllers[ctx.Reqs.Module][option.ModuleBaseControllerMapKey]; ok {
		//找不到，再查找路由指定的模块下公共控制器
		ctx.Reqs.RealController = option.ModuleBaseControllerMapKey
	} else if  obj, ok = option.RouterMap.Controllers[option.BaseModuleMapKey][option.BaseControllerMapKey]; ok {
		//再查找app下的公共控制器
		ctx.Reqs.RealModule = option.BaseModuleMapKey
		ctx.Reqs.RealController = option.BaseControllerMapKey
	} else {
		//都找不到，报错误
		AbortWithError(c, http.StatusNotFound,fmt.Sprintf("未找到对应的的控制器页面[%s.%s]", ctx.Reqs.Module, ctx.Reqs.Controller))
		return
	}

	//初始化控制器
	obj.Init(ctx)

	//判断控制器内的操作方法是否存在
	fn := reflect.ValueOf(obj).MethodByName(ToCamelCase(ctx.Reqs.Action, false))
	if fn.Kind() != reflect.Func {
		obj.AbortWithError(
			http.StatusNotFound,
			option.ErrorCallBackFunc(
				http.StatusNotFound,
				fmt.Sprintf("未找到对应的操作方法[%s.%s.%s]", ctx.Reqs.Module, ctx.Reqs.Controller, ctx.Reqs.Action),
			),
		)
		return
	}

	//调用操作方法
	fn.Call(nil)
}
