package mc

import (
	"fmt"
	"net/http"
)

//控制器接口
type IController interface {
	Init(*Context)
	AbortWithSuccess()
	AbortWithError(httpStatus int, message string, errInfo string)
	SaveLog()
}

//控制器基类
type Controller struct {
	Context  *Context
	Assign   map[string]interface{}
	Template string
	//Result *Result
}

//初始化基类
func (ctrl *Controller) Init(c *Context) {
	ctrl.Context = c
	ctrl.Template = fmt.Sprintf("%s/%s_%s.html", ctrl.Context.Reqs.RealModule, ctrl.Context.Reqs.RealController, ctrl.Context.Reqs.Action)
	ctrl.Assign = map[string]interface{}{
		"page_size": option.PageSize,
		"context":   ctrl.Context,
	}
	ctrl.SaveLog()
}

//保存日志
func (ctrl *Controller) SaveLog() {

}

func (ctrl *Controller) AbortWithSuccess() {
	ctrl.Context.HTML(http.StatusOK, ctrl.Template, ctrl.Assign)
}

func (ctrl *Controller) AbortWithError(httpStatus int, message string, errInfo string) {
	if ctrl.Assign == nil {
		ctrl.Assign = map[string]interface{}{}
	}
	ctrl.Assign["message"] = message
	ctrl.Assign["err_info"] = errInfo
	ctrl.Assign["http_status"] = httpStatus
	ctrl.Context.HTML(httpStatus, option.ErrorTemplate, ctrl.Assign)
	ctrl.Context.Abort()
}
