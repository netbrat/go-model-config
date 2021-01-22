package home

import (
    "github.com/netbrat/mc"
)

type IndexController struct {
    mc.Controller
}

func (ctrl *IndexController) IndexAct(){
    ctrl.AbortWithSuccess(mc.Result{})
}

func(ctrl *IndexController) HomeAct(){
    ctrl.AbortWithSuccess(mc.Result{})
}
