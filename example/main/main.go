package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/netbrat/go-model-config"
	"github.com/netbrat/go-model-config/example"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"time"
)


func main(){

	gin.SetMode("debug")

	option := mc.Default()

	option.RouterMap = example.RouterMap
	option.NotAuthRedirect = "/home/index/index"
	option.ErrorTemplate = "base/error.html"
	option.BaseControllerMapKey = "custom"
	option.ConfigsFilePath = "./example/model_configs/"


	_ = mc.AppendDB("default", getDB())



	m, err:= mc.NewConfigModel("sys_role")
	if err != nil{
		panic(err)
	}
	if data, err := m.GetKvs(&mc.KvsSearchOption{ExtraFields:[]string{"memo"}}); err != nil{
		panic(err)
	}else{
		fmt.Println(data)
	}
	if data, footer, total, err := m.Find(&mc.SearchOption{ExtraFields:[]string{"id"}}); err != nil {
		panic(err)
	}else{
		fmt.Println(data)
		fmt.Println(total)
		fmt.Println(footer)
	}



	//if conf,err := mc.GetConfig("sys_role"); err!=nil{
	//	panic(err)
	//}else{
	//	jsonBytes, err := json.Marshal(conf)
	//	if err != nil {
	//		panic(err)
	//	}
	//	println(string(jsonBytes))
	//	return
	//}


	return



	//r := gin.Default()
	//r.Use(Recover)
	//r.NoMethod(mc.HandlerAdapt)
	//r.NoRoute(mc.HandlerAdapt)
	//
	//r.StaticFS("/static", http.Dir("./example/static"))
	////加载模版
	//r.LoadHTMLGlob("./example/templates/admin/**/*")
	//r.GET("/",func(c *gin.Context){c.Redirect(302,"/home/index/index")})
	//
	//err = r.Run(fmt.Sprintf("%s:%d", "0.0.0.0", 8111))
	//if err != nil {
	//	panic(err)
	//}
}


func Recover(c *gin.Context){
	defer func() {
		if r := recover(); r != nil {
			mc.AbortWithError(c, http.StatusInternalServerError, fmt.Sprintf("%s",r))
		}
	}()
	c.Next()
}



func getDB() *gorm.DB{

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 0,   // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,         // 禁用彩色打印
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
	////if config.Config.Debug{
	////	db.LogMode(true)
	////}
	//dbMap[connName] = db
	//return dbMap[connName]
}