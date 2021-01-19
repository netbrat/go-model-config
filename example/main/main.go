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


func main(){


	gin.SetMode("debug")

	r := gin.Default()
	r.Use(Recover)


	r.StaticFS("/static", http.Dir("./static"))
	//加载模版
	r.SetFuncMap(mc.TemplateFuncMap)
	r.LoadHTMLGlob("./templates/**/*.html")

	r.NoMethod(mc.HandlerAdapt)
	r.NoRoute(mc.HandlerAdapt)


	option := mc.Default(r) //一定要在加载模版之后，否则默认widget不会加载
	option.ErrorTemplate = "base/error.html"
	option.ModelConfigsFilePath = "./mconfigs/"
	option.Request.PageSizeName = "limit"
	option.Request.PageName = "page"

	option.Router.ControllerMap = controller.ControllerMap
	option.Router.NotAuthActions = controller.NotAuthActions
	option.Router.NotAuthRedirect = "/home/index/index"

	option.Response.FooterName = "totalRow"
	option.Response.TotalName = "count"
	option.Response.SuccessCodeValue = "0"

    option.Auth.RowAuthModels = controller.RowAuthModels
	option.Auth.GetAuthFunc = func() *mc.Auth{ return &mc.Auth{IsSuper:true}}

	_ = mc.AppendDB("default", getDB())

	r.GET("/",func(c *gin.Context){c.Redirect(302,"/home/index/index")})

	err := r.Run(fmt.Sprintf("%s:%d", "0.0.0.0", 8111))
	if err != nil {
		panic(err)
	}
}


func Recover(c *gin.Context){
	defer func() {
		if r := recover(); r != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println(r)
			fmt.Println(string(buf[:n]))
			mc.AbortWithErrorService(c, mc.Result{HttpStatus:500, Code:"500", Message:fmt.Sprintf("%s",r)})
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
	//
	////分析一个连接配置信息
	//masters := make([]gorm.Dialector, 0)
	//salves := make([]gorm.Dialector,0)
	//var defDialector gorm.Dialector
	//for _ , v := range dbConnConfig.DbConfigs {
	//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
	//		v.User, v.Password, v.Host, v.Port, v.DbName, v.Charset,
	//	)
	//	dialector := mysql.New(mysql.Config{DSN: dsn, DisableDatetimePrecision: true})
	//	if !v.Salves && defDialector == nil{
	//		defDialector = dialector
	//		continue
	//	}
	//	if v.Salves {
	//		salves = append(salves, dialector)
	//	}else{
	//		masters = append(masters, dialector)
	//	}
	//}
	//if defDialector == nil {
	//	panic(fmt.Sprintf("数据库连接未设置主库(%s)", connName))
	//}
	//
	//// 初始化一个连接
	//db, err := gorm.Open(defDialector, &gorm.Config{Logger:newLogger})
	//if err != nil {
	//	panic(fmt.Sprintf("数据库连接失败(%s):%s", connName, err))
	//}
	//// 设置主从
	//if err = db.Use(dbresolver.Register(dbresolver.Config{
	//	Sources:  masters,
	//	Replicas: salves,
	//	Policy:   dbresolver.RandomPolicy{},
	//})); err != nil {
	//	panic(fmt.Sprintf("数据库连接失败(%s):%s", connName, err))
	//}
	//
	//// 设置连接选项
	//if sqlDB, err := db.DB(); err == nil && sqlDB != nil{
	//	sqlDB.SetMaxIdleConns(dbConnConfig.MaxIdleConns)
	//	sqlDB.SetMaxOpenConns(dbConnConfig.MaxOpenConns)
	//	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(dbConnConfig.ConnMaxLifetime))
	//	sqlDB.SetConnMaxIdleTime(time.Minute * time.Duration(dbConnConfig.ConnMaxIdleTime))
	//}else{
	//	panic(fmt.Sprintf("数据库连接失败(%s):%s", connName, err))
	//}
	//
	//var data []map[string]interface{}
	//db.Table("sys_role").Find(&data)
	//fmt.Println(data)
	//
	//
	//
	////if configs.Config.Debug{
	////	db.LogMode(true)
	////}
	//dbMap[connName] = db
	//return dbMap[connName]
}