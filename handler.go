package mc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

//未初始化控制器之前的错误处理
func AbortWithError(c *gin.Context, httpStatus int, code int, err error){
	//初始化上下文
	obj := &Controller{}
	obj.Initialize(NewContext(c))
	obj.AbortWithError(httpStatus, code, err)
}


//所有入口适配器
func HandlerAdapt(c *gin.Context) {
	//初始化上下文
	ctx := NewContext(c)
	//判断是否需要登录验证
	actionString := fmt.Sprintf("%s.%s.%s", ctx.ModelName, ctx.ControllerName, ctx.ActionName)
	if !inArray(actionString, option.NotAuthActions) {
		//验证是否登录
	}

	//从注册表中查询路由指定的控制器
	var obj IController
	ok := false
	//先查找路由指定的模块及控制器
	if obj, ok = option.ControllerMap[ctx.ModuleName][ctx.ControllerName]; ok {
		//找到不作处理
	} else if obj, ok = option.ControllerMap[ctx.ModuleName][option.ModuleBaseControllerMapKey]; ok {
		//找不到，再查找路由指定的模块下公共控制器
		ctx.RealControllerName = option.ModuleBaseControllerMapKey
	} else if  obj, ok = option.ControllerMap[option.BaseModuleMapKey][option.BaseControllerMapKey]; ok {
		//再查找app下的公共控制器
		ctx.RealModuleName = option.BaseModuleMapKey
		ctx.RealControllerName = option.BaseControllerMapKey
	} else {
		//都找不到，报错误
		AbortWithError(c, http.StatusNotFound, http.StatusNotFound, fmt.Errorf(fmt.Sprintf("未找到对应的的控制器页面[%s.%s]", ctx.ModuleName, ctx.ControllerName)))
		return
	}

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
			fmt.Println("ERR")
			obj.AbortWithError(http.StatusNotFound, http.StatusNotFound, fmt.Errorf(fmt.Sprintf("未找到对应的操作方法[%s.%s.%s]", ctx.ModuleName, ctx.ControllerName, ctx.ActionName)))
			return
		}
	}

	//调用操作方法
	fn.Call(nil)
}
