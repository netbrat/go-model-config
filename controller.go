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
	InitializeBefore()                                 //初始化Controller前回调
	Initialize(c *Context, childCtrl IController)      //初始化
	InitializeAfter()                                  //初始化Controller后回调
	ModelInitializeBefore()                            //初始化Model前回调
	ModelInitializeAfter()                             //初始化Model后回调
	ModelListUIRenderBefore()                          //列表界面渲染前回调
	ModelFindBefore(qo *QueryOption)                   //查询数据前回调
	ModelFindAfter(result *Result)                     //查询数据后回调
	ModelEditTakeBefore(qo *QueryOption)               //编辑查询数据前回调
	ModelEditTakeAfter(rowData map[string]interface{}) //编辑查询数据后回调
	ModelSaveBefore(data map[string]interface{})       //保存数据前回调
	ModelSaveAfter(result *Result)                     //保存数据后回调
	ModelDelBefore(ids []string)                       //删除数据前回调
	ModelDelAfter(result *Result)                      //删除数据后回调
	AbortWithSuccess(result *Result)                   //成功输出
	AbortWithError(err interface{})                    //错误输出
	AbortWithMessage(result *Result)                   //消息输出 （使用消息模版）
	AbortWithHtml(result *Result)                      //html输出 （使用HTML渲染式）
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

//结果体，当做为错误panic
type Result struct {
	HttpStatus  int                    //响应代码
	Code        string                 //消息代码
	Message     string                 //消息
	RedirectUrl string                 //跳转地址
	Data        interface{}            //数据体
	ExtraData   map[string]interface{} //附加的数据
	RenderType  RenderType             //渲染方式
	IsInfo		bool					//是否是普通提示性消息（使用错误输出时有效）
}

//给分配器添加内容
func (assign * Assign) Append(key string, value interface{}) (r *Assign){
	assign.Result[key] = value
	return assign
}

//实现结果错误接口
func (result *Result) Error() string {
	return result.Message
}

//给结果集上添加内容
func (result *Result) Append(key string, value interface{}) (r *Result){
	result.ExtraData[key] = value
	return result
}

//初始化Controller前回调
func (ctrl *Controller) InitializeBefore(){}

//初始化Controller后回调
func (ctrl *Controller) InitializeAfter(){}

//初始化Model前回调
func (ctrl *Controller) ModelInitializeBefore() {}

//初始化Model后回调
func (ctrl *Controller) ModelInitializeAfter() {}

//列表界面渲染前回调
func(ctrl *Controller) ModelListUIRenderBefore(){
	return
}
//查询数据前回调
func(ctrl *Controller) ModelFindBefore(qo *QueryOption){
	return
}

//查询数据后回调
func(ctrl *Controller) ModelFindAfter(result *Result){
	return
}

//编辑查询数据前回调
func(ctrl *Controller) ModelEditTakeBefore(qo *QueryOption){
	return
}
//编辑查询数据后回调
func(ctrl *Controller) ModelEditTakeAfter(rowData map[string]interface{}){
	return
}

//保存数据前回调
func(ctrl *Controller) ModelSaveBefore(data map[string]interface{}){
	return
}

//保存数据后回调
func(ctrl *Controller) ModelSaveAfter(result *Result){
	return
}

//删除数据前回调
func(ctrl *Controller) ModelDelBefore(ids []string){
	return
}

//删除数据后回调
func(ctrl *Controller) ModelDelAfter(result *Result){
	return
}

//初始化基类
func (ctrl *Controller) Initialize(c *Context, childCtrl IController) {
	ctrl.ChildController = childCtrl
	//初始化之前
	ctrl.ChildController.InitializeBefore()
	//开始初始化
	ctrl.Context = c
	ctrl.Template = "" // fmt.Sprintf("%s/%s_%s.html", ctrl.Context.RealModuleName, ctrl.Context.RealControllerName, ctrl.Context.ActionName)
	ctrl.Assign = &Assign{
		Context: c,
		Model: ctrl.Model,
		Result: map[string]interface{}{
			option.Response.CodeName: option.Response.SuccessCodeValue,
			option.Response.MessageName: "",
			option.Response.DataName: nil,
		},
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
	}else if e, ok := err.(error); ok {
		result = &Result{HttpStatus: http.StatusInternalServerError, Code: cast.ToString(http.StatusInternalServerError), Message: e.Error()}
	}else {
		result = &Result{HttpStatus: http.StatusInternalServerError, Code: cast.ToString(http.StatusInternalServerError), Message: "系统异常，请稍后重试!"}
	}
	//消息代码
	if result.Code == "" {
		result.Code = option.Response.FailCodeValue
	}
	if result.IsInfo { //使用普通消息输出
		ctrl.AbortWithMessage(result)
		return
	}
	ctrl.Template = option.Response.ErrorTemplate
	ctrl.AbortWithSuccess(result)
}

//消息输出 （使用消息模版）
func (ctrl *Controller) AbortWithMessage(result *Result) {
	fmt.Println("message", ctrl.Context)
	ctrl.Template = option.Response.MessageTemplate
	ctrl.AbortWithSuccess(result)
}

//html输出 （使用HTML渲染式）
func (ctrl *Controller) AbortWithHtml(result *Result){
	fmt.Println("html", ctrl.Context)
	if result == nil{
		result = &Result{}
	}
	result.RenderType = RenderTypeHTML
	ctrl.AbortWithSuccess(result)
}





