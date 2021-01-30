package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/netbrat/mc"
	"github.com/netbrat/mc/example/controller"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)


func main() {
	gin.SetMode("debug")

	r := gin.Default()
	//r.Use(Recover)

	r.StaticFS("/static", http.Dir("./static"))
	//加载模版
	r.SetFuncMap(mc.TemplateFuncMap)
	r.LoadHTMLGlob("./templates/**/*.html")

	r.NoMethod(mc.HandlerAdapt)
	r.NoRoute(mc.HandlerAdapt)

	option := mc.Default(r) //一定要在加载模版之后，否则默认widget不会加载
	option.Response.ErrorTemplate = "base/error.html"
	option.Response.MessageTemplate = "base/message.html"
	option.ModelConfigsFilePath = "./mconfigs/"
	option.Request.PageSizeName = "limit"
	option.Request.PageName = "page"

	option.Router.ControllerMap = controller.ControllerMap

	option.Response.FootName = "totalRow"
	option.Response.TotalName = "count"
	option.Response.SuccessCodeValue = "0"

	//option.ModelAuth.RowAuthModels = controller.RowAuthModels
	//option.ModelAuth.GetModelAuthCallback = func() *mc.ModelAuth {
	//	return mc.NewModelAuth(true, nil, nil)
	//}

	_ = mc.AppendDB("default", getDB())

	r.GET("/", func(c *gin.Context) { c.Redirect(302, "/home/index/index") })

	err := r.Run(fmt.Sprintf("%s:%d", "0.0.0.0", 8111))
	if err != nil {
		panic(err)
	}
}


func Recover(c *gin.Context){
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println(string(buf[:n]))
			mc.AbortWithError(c, err)
		}
	}()
	c.Next()
}



func getDB() *gorm.DB{

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 0,           // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	dsn := "root:123456@tcp(127.0.0.1:3306)/mc_test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(
		mysql.New(mysql.Config{DSN:dsn, DisableDatetimePrecision:true}),
		&gorm.Config{Logger:newLogger},
	)
	if err != nil {
		panic(fmt.Sprintf("数据库连接错误"))
	}

	if sqlDB, err := db.DB(); err == nil && sqlDB != nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(20))
		sqlDB.SetConnMaxIdleTime(time.Minute * time.Duration(20))
	}else {
		panic(fmt.Sprintf("数据库连接失败"))
	}



	return db
}



