package example

import (
	mc "github.com/netbrat/go-model-config"
	"github.com/netbrat/go-model-config/example/home/controllers"
)

//路由控制器注册表
var RouterMap = mc.RouterMap{
	//控制器注册表：格式 map[module][controller]
	Controllers:    map[string]map[string]mc.IController{
		"base": {
			"custom": &controllers.IndexController{},
		},
		"home": {
			"index": &controllers.IndexController{},
		},
	},
	//不需要登录验证的操作方法
	NotAuthActions: []string{
		"home.public.login",
	},
}