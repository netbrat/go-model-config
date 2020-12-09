package mc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type RouterMap struct {
	Controllers    map[string]map[string]IController
	NotAuthActions []string
}

//所有入口适配器
func HandlerAdapt(c *gin.Context) {
	//初始化上下文
	ctx := &Context{Context: c}
	ctx.Init()
	//判断是否需要登录验证
	actionString := fmt.Sprintf("%s.%s.%s", ctx.Reqs.Model, ctx.Reqs.Controller, ctx.Reqs.Action)
	if !inArray(actionString, option.RouterMap.NotAuthActions) {
		//验证是否登录
	}

	//从注册表中查询路由指定的控制器
	var obj IController
	//var obj interface{}
	ok := false
	//先查找路由指定的模块及控制器
	if obj, ok = option.RouterMap.Controllers[ctx.Reqs.Module][ctx.Reqs.Controller]; !ok {
		//再查找路由指定的模块下公共控制器
		if obj, ok = option.RouterMap.Controllers[ctx.Reqs.Module]["base"]; !ok {
			//再查找app下的公共控制器
			if obj, ok = option.RouterMap.Controllers["base"]["custom"]; !ok {
				obj = &Controller{}
				obj.Init(&Context{Context: c})
				obj.AbortWithError(http.StatusNotFound, fmt.Sprintf("未找到对应的的控制器页面[%s.%s]", ctx.Reqs.Module, ctx.Reqs.Controller), "")
				return
			} else {
				ctx.Reqs.RealModule = "base"
				ctx.Reqs.RealController = "custom"
			}
		} else {
			ctx.Reqs.RealController = "base"
		}
	}

	//初始化控制器
	obj.Init(ctx)

	//判断控制器内的操作方法是否存在
	fn := reflect.ValueOf(obj).MethodByName(ToCamelCase(ctx.Reqs.Action, false))
	if fn.Kind() != reflect.Func {
		obj.AbortWithError(http.StatusNotFound, fmt.Sprintf("未找到对应的操作方法[%s.%s.%s]", ctx.Reqs.Module, ctx.Reqs.Controller, ctx.Reqs.Action), "")
		return
	}

	//调用操作方法
	fn.Call(nil)
}
