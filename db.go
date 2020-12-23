package mc

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type DbConfig struct {
	User         string		`json:"user"`
	Password     string		`json:"password"`
	Host         string		`json:"host"`
	Port         int		`json:"port"`
	DbName       string		`json:"db_name"`
	Charset      string		`json:"charset"`
	MaxOpenConns int		`json:"max_open_conns"`
	MaxIdleConns int		`json:"max_idle_conns"`
}

type DbOption struct {
	DB				*gorm.DB
	ConnName 		string
	DbName			string
	Table			string
	Pk				string
	AutoIncrement 	bool
	UniqueFields	[]string
}

func (do *DbOption) Set(connName string, dbName string, table string, pk string, autoIncrement bool, uniqueFields []string) {
	do.ConnName = connName
	do.Table = table
	do.DB = GetDB(connName)
	if dbName == ""{
		dbName = GetDbNameByConfig(connName)
	}
	do.DbName = dbName
	do.AutoIncrement = autoIncrement
	do.UniqueFields = uniqueFields
}

//数据库连接配置map
var dbConfigMap = map[string]DbConfig{}

//数据库连接池
var dbMap = map[string]*gorm.DB{}

//根据数据库配置名，获取数据库名称
func GetDbNameByConfig(connName string) string {
	dbConfig, ok := dbConfigMap[connName]
	if !ok{
		panic(fmt.Sprintf("数据库连接配置项不存在[%s]", connName))
	}
	return dbConfig.DbName
}

//添加数据库连接配置
func AppendDbConfig(connName string , dbConfig DbConfig){
	dbConfigMap[connName] = dbConfig
}

//获取一个数据库连接对象
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
