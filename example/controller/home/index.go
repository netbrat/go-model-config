package home

import (
    "github.com/netbrat/mc"
)

type IndexController struct {
    mc.Controller
}

func (ctrl *IndexController) Index(){
    ctrl.AbortWithSuccess(mc.Result{})
}

func(ctrl *IndexController) Home(){
    ctrl.AbortWithSuccess(mc.Result{})
}
