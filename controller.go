package mc

import (
	"fmt"
	"net/http"
)

//控制器接口
type IController interface {
	Init(*Context)
	AbortWithSuccess(result interface{})
	AbortWithError(httpStatus int, result interface{})
	SaveLog()
}


//输出内容
type Assign struct {
	HttpStatus	int
	Context		*Context
	Result		interface{}
}


//控制器基类
type Controller struct {
	Context  	*Context
	Assign   	*Assign
	Template 	string
}

//初始化基类
func (ctrl *Controller) Init(c *Context) {
	ctrl.Context = c
	ctrl.Template = fmt.Sprintf("%s/%s_%s.html", ctrl.Context.Reqs.RealModule, ctrl.Context.Reqs.RealController, ctrl.Context.Reqs.Action)
	ctrl.Assign = &Assign{Context: c}
	ctrl.SaveLog()
}

//保存日志
func (ctrl *Controller) SaveLog() {

}

//成功输出
func (ctrl *Controller) AbortWithSuccess(result interface{}) {
	if result != nil {
		ctrl.Assign.Result = result
	}else{
		ctrl.Assign.Result = map[string]interface{}{}
	}
	ctrl.Assign.HttpStatus = http.StatusOK
	ctrl.Context.HTML(http.StatusOK, ctrl.Template, ctrl.Assign)
	ctrl.Context.Abort()
}

//错误输出
func (ctrl *Controller) AbortWithError(httpStatus int, result interface{}) {
	if result != nil {
		ctrl.Assign.Result = result
	}else{
		ctrl.Assign.Result = map[string]interface{}{}
	}
	ctrl.Assign.HttpStatus = httpStatus
	ctrl.Context.HTML(httpStatus, option.ErrorTemplate, ctrl.Assign)
	ctrl.Context.Abort()
}
