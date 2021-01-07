package mc

import (
	"fmt"
	"net/http"
)

//控制器接口
type IController interface {
	Initialize(*Context)
	AbortWithSuccess(result map[string]interface{})
	AbortWithError(httpStatus int, code int, err error)
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

//结果
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Footer  interface{} `json:"footer"`
	Other	interface{} `json:"other"`
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

//成功输出
func (ctrl *Controller) AbortWithSuccess(result map[string]interface{}) {
	//result := option.SuccessCallBackFunc(ctrl.Context, r)
	ctrl.Assign.Result = result
	ctrl.Context.Render(http.StatusOK, ctrl.Template, ctrl.Assign)
	ctrl.Context.Abort()
}

//错误输出
func (ctrl *Controller) AbortWithError(httpStatus int, code int, err error) {
	ctrl.Assign.Result = map[string]interface{}{
		"code":    code,
		"message": err.Error(),
	}
	ctrl.Context.Render(httpStatus, option.ErrorTemplate, ctrl.Assign)
	ctrl.Context.Abort()
}



