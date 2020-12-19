package controllers

import (
	mc "github.com/netbrat/go-model-config"
)

type IndexController struct {
	mc.Controller
}

func (ctrl *IndexController) Index(){
	ctrl.AbortWithSuccess(nil)
}

func(ctrl *IndexController) Biz(){
	ctrl.AbortWithSuccess(nil)
}