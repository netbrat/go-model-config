package controller

import (
    "github.com/netbrat/mc"
    "github.com/netbrat/mc/example/controller/base"
    "github.com/netbrat/mc/example/controller/home"
)

// 控制器注册表：格式 map[module][controller]
var ControllerMap = map[string]map[string]mc.IController{
    "base": {
        "base": &base.ModelController{},
    },
    "home": {
        "index": &home.IndexController{},
    },
}

// 无需登录的操作方法
var NotAuthActions = []string{
    "home.public.login",
}
