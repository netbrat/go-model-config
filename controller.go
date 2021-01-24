package mc

import (
	"fmt"
	"github.com/spf13/cast"
	"net/http"
	"net/url"
	"strings"
)

//控制器接口
type IController interface {
	InitializeBefore()
	InitializeAfter()
	ModelInitializeBefore()
	ModelInitializeAfter()
	Initialize(c *Context, childCtrl IController)
	AbortWithSuccess(result *Result)
	AbortWithError(err interface{})
}


//控制器基类
type Controller struct {
	Context  *Context
	Model *Model
	Template string
	ChildController IController
	Assign   *Assign
}


//内容分配
type Assign struct {
	Context *Context
	Model   *Model
	Result  map[string]interface{}
}

//初始化Controller之前
func (ctrl *Controller) InitializeBefore(){}

//初始化Controller之后
func (ctrl *Controller) InitializeAfter(){}

//初始化Model之前
func (ctrl *Controller) ModelInitializeBefore() {}

//初始化Model之后
func (ctrl *Controller) ModelInitializeAfter() {}


//初始化基类
func (ctrl *Controller) Initialize(c *Context, childCtrl IController) {
	ctrl.ChildController = childCtrl
	//初始化之前
	ctrl.ChildController.InitializeBefore()
	//开始初始化
	ctrl.Context = c
	ctrl.Template = fmt.Sprintf("%s/%s_%s.html", ctrl.Context.RealModuleName, ctrl.Context.RealControllerName, ctrl.Context.ActionName)
	ctrl.Assign = &Assign{
		Context: c,
		Model: ctrl.Model,
		Result: make(map[string]interface{}),
	}
	//初始化之后
	ctrl.ChildController.InitializeAfter()
}

//UrlValue 转换成 RequestValue
func (ctrl *Controller) UrlValueToRequestValue(values url.Values) (requestValue map[string]interface{}) {
	requestValue = make(map[string]interface{})
	for key, value := range values {
		key := strings.ReplaceAll(key, "[]", "")
		if len(value) <= 1 {
			requestValue[key], _ = url.QueryUnescape(value[0])
		} else {
			for i, _ := range value {
				value[i], _ = url.QueryUnescape(value[i])
			}
			requestValue[key] = value
		}
	}
	return
}

//使用默认的结果格式输出
func (ctrl *Controller) AbortWithSuccess(result *Result) {
	if result == nil{
		result = &Result{}
	}
	//响应代码
	httpStatus := http.StatusOK
	if result.HttpStatus != 0 {
		httpStatus = result.HttpStatus
	}
	//消息代码
	if result.Code != "" {
		ctrl.Assign.Result[option.Response.CodeName] = result.Code
	}
	if cast.ToString(ctrl.Assign.Result[option.Response.CodeName]) == "" {
		ctrl.Assign.Result[option.Response.CodeName] = option.Response.SuccessCodeValue
	}

	//消息
	if result.Message != "" {
		ctrl.Assign.Result[option.Response.MessageName] = result.Message
	}

	//数据
	if result.Data != nil {
		ctrl.Assign.Result[option.Response.DataName] = result.Data
	}

	//扩展数据
	if result.ExtraData != nil {
		for key, value := range result.ExtraData {
			if key == option.Response.MessageName || key == option.Response.CodeName || key == option.Response.DataName {
				continue
			}
			ctrl.Assign.Result[key] = value
		}
	}
	ctrl.Context.Render(result.RenderType, httpStatus, ctrl.Template, ctrl.Assign)
	ctrl.Context.Abort()
}

//错误输出（标准错误, 页面状态还是200）
func (ctrl *Controller) AbortWithError(err interface{}) {
	var result *Result
	if err == nil {
		result = &Result{}
	}else if e, ok := err.(*Result); ok {
		result = e
	}else {
		result = &Result{HttpStatus: http.StatusInternalServerError, Code: cast.ToString(http.StatusInternalServerError), Message: "系统异常，请稍后重试!"}
	}
	//消息代码
	if result.Code == "" {
		result.Code = option.Response.FailCodeValue
	}
	ctrl.Template = option.ErrorTemplate
	ctrl.AbortWithSuccess(result)
}






