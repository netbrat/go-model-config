package mc

import "fmt"

type ModelController struct {
    Controller
    Model *Model
}


func (ctrl *ModelController) Initialize(c *Context) {
    if c.ModelName == "" {
        c.ModelName = fmt.Sprintf("%s_%s", c.ModuleName,c.ControllerName)
    }
    ctrl.Controller.Initialize(c)
    ctrl.Model = NewModel(ctrl.Context.ModelName)
    ctrl.Assign.Model = ctrl.Model
}


func (ctrl *ModelController) Index(){
    if ctrl.Context.Request.Method=="GET"{
        //r := &Result{}
        //ctrl.AbortWithSuccess(r)
        ctrl.Model.CreateSearchItems(nil)
        ctrl.AbortWithSuccess(nil)

    }else{

    }
}


func (ctrl *ModelController) Add(){
    if ctrl.Context.Request.Method=="GET"{
        ctrl.Model.CreateEditItems(nil)
        ctrl.AbortWithSuccess(nil)
    }else{
        ctrl.AbortWithError(200,200,fmt.Errorf("abc"))
        fmt.Println(ctrl.Context.Request.Form)
        ctrl.Context.RenderType = RenderTypeString
        result := map[string]interface{}{
            "search": ctrl.Context.Request.Form,
        }
        ctrl.AbortWithSuccess(result)
    }
}

func (ctrl *ModelController) Edit(){
    if ctrl.Context.Request.Method=="GET"{

    }else{

    }
}