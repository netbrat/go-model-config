package mc

import "fmt"

type ModelController struct {
    Controller
    Model *Model
}


func (ctrl *ModelController) Initialize(c *Context) {
    ctrl.Controller.Initialize(c)
    ctrl.Model = NewModel(ctrl.Context.ModelName)
    ctrl.Assign.Model = ctrl.Model
}


func (ctrl *ModelController) Index(){
    if ctrl.Context.Request.Method=="GET"{
        //r := &Result{}
        //ctrl.AbortWithSuccess(r)
        ctrl.Model.CreateSearchItems(nil)
        fmt.Println(ctrl.Model.SearchItems)
        ctrl.AbortWithSuccess(nil)

    }else{

    }
}
