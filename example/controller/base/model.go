package base

import "github.com/netbrat/mc"

type ModelController struct {
    mc.ModelController
}


func (ctrl *ModelController) Initialize(c *mc.Context){
    c.ModelName = "sys_role"
    ctrl.ModelController.Initialize(c)
}
