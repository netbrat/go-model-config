package mc

import (
	"fmt"
	"net/url"
	"strings"
)

//控制器接口
type IController interface {
	Initialize(*Context)
	AbortWithSuccess(result Result)
	AbortWithError(result Result)
	SaveLog()
}


//控制器基类
type Controller struct {
	Context  	*Context
	Template 	string
	Auth		*Auth
	Assign		*Assign
	LogIgnoreActions []string	//不保存的操作方法，如["*"]表示该控制器下的所有方法全部不保存，默认["index","export"]
}


type Assign struct {
	Context *Context
	Model   *Model
	Result  map[string]interface{}
}


type Result struct {
	HttpStatus int                    //响应代码
	Code       string                 //消息代码
	Message    string                 //消息
	Data       interface{}            //数据体
	ExtraData  map[string]interface{} //附加的数据
	RenderType RenderType             //渲染方式
}


//初始化基类
func (ctrl *Controller) Initialize(c *Context) {
	ctrl.Context = c
	ctrl.Template = fmt.Sprintf("%s/%s_%s.html", ctrl.Context.RealModuleName, ctrl.Context.RealControllerName, ctrl.Context.ActionName)
	if ctrl.LogIgnoreActions == nil || len(ctrl.LogIgnoreActions)<=0 {
		ctrl.LogIgnoreActions = []string{"index","export"}
	}
	ctrl.Assign = &Assign{Context: ctrl.Context}
	ctrl.SaveLog()
}

//保存日志
func (ctrl *Controller) SaveLog() {

}

func (ctrl *Controller) UrlValueToRequestValue(values url.Values) (searchValue map[string]interface{}){
	searchValue = make(map[string]interface{})
	for key, value := range values{
		key := strings.ReplaceAll(key, "[]","")
		if len(value)<=1 {
			searchValue[key], _ = url.QueryUnescape(value[0])
		}else{
			for i, _ := range value {
				value[i], _ = url.QueryUnescape(value[i])
			}
			searchValue[key] = value
		}
	}
	return
}



//使用默认的结果格式输出
func (ctrl *Controller) AbortWithSuccess(result Result){
	//响应代码
	httpStatus := 200
	if result.HttpStatus != 0 {
		httpStatus = result.HttpStatus
	}
	//结果
	newResult := map[string]interface{}{
		option.Response.CodeName: result.Code,
		option.Response.MessageName: result.Message,
		option.Response.DataName: result.Data,
	}
	//消息代码
	if result.Code == "" {
		newResult[option.Response.CodeName] = option.Response.SuccessCodeValue
	}
	//扩展数据
	if result.ExtraData != nil{
		for key, value := range result.ExtraData{
			if key == option.Response.MessageName || key == option.Response.CodeName || key == option.Response.DataName {
				continue
			}
			newResult[key] = value
		}
	}
	ctrl.Assign.Result = newResult
	ctrl.Context.Render(result.RenderType, httpStatus, ctrl.Template, ctrl.Assign)
	ctrl.Context.Abort()
}

//错误输出（标准错误, 页面状态还是200）
func (ctrl *Controller) AbortWithError(result Result) {
	//消息代码
	if result.Code == "" {
		result.Code = option.Response.FailCodeValue
	}
	ctrl.Template = option.ErrorTemplate
	ctrl.AbortWithSuccess(result)
}






