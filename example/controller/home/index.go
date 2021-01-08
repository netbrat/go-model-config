package home

import "github.com/netbrat/mc"

type IndexController struct {
    mc.Controller
}

func (ctrl *IndexController) Index(){
    ctrl.AbortWithSuccess(nil)
}

func(ctrl *IndexController) Manage(){
    ctrl.AbortWithSuccess(nil)
}
