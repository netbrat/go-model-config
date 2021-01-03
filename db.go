package mc

import (
	"fmt"
	"gorm.io/gorm"
)

//数据库连接池
var dbMap = map[string]*gorm.DB{}


// 添加一个数据库连接对象
func AppendDB(connName string, db *gorm.DB) (err error) {
	if connName == "" || db == nil{
		err = fmt.Errorf("不是有效的数据库连接名和数据库连接对象")
	}else {
		dbMap[connName] = db
	}
	return
}

// 获取一个数据库连接对象
func GetDB(connName string) (db *gorm.DB, err error){
	ok := false
	if db, ok = dbMap[connName]; ok {
		db = db.Table("")
	}else{
		err = fmt.Errorf("数据库连接项不存在[%s]", connName)
	}
	return
}
