package controllers

import mc "github.com/netbrat/go-model-config"

type TestController struct {
    mc.ModelController
}

func (ctrl *TestController) Initialize(c *mc.Context){
    c.ModelName = "sys_role"
    ctrl.ModelController.Initialize(c)

}

