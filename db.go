package mc

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type DbConfig struct {
	User         string
	Password     string
	Host         string
	Port         int
	DbName       string
	Charset      string
	MaxOpenConns int
	MaxIdleConns int
}

var dbMap = map[string]*gorm.DB{}
var dbConfigMap = map[string]DbConfig{}

func GetDB(connName string) *gorm.DB {
	dbConfig, ok := dbConfigMap[connName]
	if !ok {
		panic(fmt.Sprintf("数据库连接配置项不存在[%s]", connName))
	}
	//判断连接Map中是否存在,存在则直接使用
	if db, ok := dbMap[connName]; ok {
		return db
	}
	//不存在则创建新连接
	mysqlDbDns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName,
		dbConfig.Charset,
	)
	db, err := gorm.Open("mysql", mysqlDbDns)
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败('%s):(%s)", mysqlDbDns, err))
	}
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.DB().SetMaxOpenConns(dbConfig.MaxIdleConns)
	//if config.Config.Debug{
	//	db.LogMode(true)
	//}
	dbMap[connName] = db
	return dbMap[connName]
}
