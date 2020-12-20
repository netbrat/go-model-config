package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/netbrat/go-model-config"
	"github.com/netbrat/go-model-config/example"
	"net/http"
)

func main(){

	gin.SetMode("debug")

	option := mc.Default()

	option.RouterMap = example.RouterMap
	option.NotAuthRedirect = "/home/index/index"
	option.ErrorTemplate = "base/error.html"
	option.BaseControllerMapKey = "custom"
	option.ConfigsFilePath = "./example/model_configs/"

	if conf,err := mc.GetConfig("sys_role"); err!=nil{
		panic(err)
	}else{
		jsonBytes, err := json.Marshal(conf)
		if err != nil {
			panic(err)
		}
		println(string(jsonBytes))
		return
	}



	r := gin.Default()
	r.Use(Recover)
	r.NoMethod(mc.HandlerAdapt)
	r.NoRoute(mc.HandlerAdapt)

	r.StaticFS("/static", http.Dir("./example/static"))
	//加载模版
	r.LoadHTMLGlob("./example/templates/admin/**/*")
	r.GET("/",func(c *gin.Context){c.Redirect(302,"/home/index/index")})

	err := r.Run(fmt.Sprintf("%s:%d", "0.0.0.0", 8111))
	if err != nil {
		panic(err)
	}
}


func Recover(c *gin.Context){
	defer func() {
		if r := recover(); r != nil {
			mc.AbortWithError(c, http.StatusInternalServerError, fmt.Sprintf("%s",r))
		}
	}()
	c.Next()
}

