package example

import (
	mc "github.com/netbrat/go-model-config"
	"github.com/netbrat/go-model-config/example/home/controllers"
)


// 控制器注册表：格式 map[module][controller]
var Controllers = map[string]map[string]mc.IController{
	"base": {
		"custom": &controllers.IndexController{},
	},
	"home": {
		"index": &controllers.IndexController{},
		"test": &controllers.TestController{},
	},
}

// 无需登录的操作方法
var NotAuthActions = []string{
	"home.public.login",
}
