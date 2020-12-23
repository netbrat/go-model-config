package mc

import (
	"fmt"
	"reflect"
)

type SearchOption struct {
	Where		string			//查询条件
	WhereValue	[]interface{}	//查询值
	Fields		[]string		//查询字段
	Page		int				//查询页码
	PageSize 	int				//分页大小
	OrderBy		string			//排序
	Join		string			//外联
	Group		string			//分组
	Alias		string			//别名
	Having		string
	NotTotal	bool			//是否不查询总记录数
}

func (so SearchOption) Offset() int{
	return (so.Page-1) * so.PageSize
}

type BaseModel struct {
	DbOption 	*DbOption
}
//设置数据库操作选项
func (bm *BaseModel) SetDbOption(connName string, dbName string, table string, pk string, autoIncrement bool, uniqueFields []string) *BaseModel{
	bm.DbOption = &DbOption{}
	bm.DbOption.Set(connName, dbName, table, pk, autoIncrement, uniqueFields)
	return bm
}

//获取单条记录
func (bm *BaseModel) First(so SearchOption)(data map[string]interface{}, err error){
	err = bm.DbOption.DB.
		Table(bm.DbOption.Table).
		Order(so.OrderBy).
		Select(so.Fields).
		Where(so.Where,so.WhereValue...).
		Joins(so.Join).
		Having(so.Having).
		First(data).Error
	return
}

//获取记录集
func (bm *BaseModel) Find(so SearchOption)(data []map[string]interface{}, total int, err error){
	db := bm.DbOption.DB.
		Table(bm.DbOption.Table).
		Order(so.OrderBy).
		Select(so.Fields).
		Where(so.Where,so.WhereValue...).
		Joins(so.Join).
		Limit(so.PageSize).
		Offset(so.Offset()).
		Having(so.Having).
		Find(data)
	if err = db.Error; err != nil{
		return
	}
	if !so.NotTotal {
		err = db.Count(&total).Error
	}
	return
}

//判断记录是否存在
func (bm *BaseModel) IsExist(data map[string]interface{}) (exist bool, err error){
	where := ""
	whereValue := make([]interface{},0)
	db := bm.DbOption.DB.Table(bm.DbOption.Table)
	for _,v := range bm.DbOption.UniqueFields {
		if where == "" {
			where += fmt.Sprintf(" AND %s = ?", v)
		}else{
			where = fmt.Sprintf("%s = ?", v)
		}
		whereValue = append(whereValue, data[v])
	}

	if !bm.DbOption.AutoIncrement{
		where = fmt.Sprintf("(%s) OR (%s = ?)", where, bm.DbOption.Pk)
		whereValue = append(whereValue, data[bm.DbOption.Pk])
	}
	total := 0
	db = db.Where(where, whereValue...).Count(&total)
	if total >0 {
		exist = true
	}
	return exist, db.Error
}

//更新记录
func (bm *BaseModel) Update(data map[string]interface{}, id interface{})(total int64, err error){
	exist := false
	if exist, err = bm.IsExist(data); err != nil{
		return
	}else if exist {
		err = fmt.Errorf("记录已存在")
		return
	}
	where := fmt.Sprintf("%s = ?", bm.DbOption.Pk)
	db := bm.DbOption.DB.Table(bm.DbOption.Table).Where(where, id).Update(data)
	return db.RowsAffected, db.Error
}

//创建记录
func (bm *BaseModel) Create(data map[string]interface{})(total int64, err error){
	exist := false
	if exist, err = bm.IsExist(data); err != nil{
		return
	}else if exist {
		err = fmt.Errorf("记录已存在")
		return
	}
	db := bm.DbOption.DB.Table(bm.DbOption.Table).Create(data)
	return db.RowsAffected, db.Error
}

//保存记录（根据pk自动分析是update 或 create）
func (bm *BaseModel) Save(data map[string]interface{})(total int64, err error){
	pk := ""
	where := map[string]interface{}{}
	if bm.DbOption.AutoIncrement { //pk自增表
		pk = bm.DbOption.Pk
	}else{
		pk = "__" + bm.DbOption.Pk
		where[bm.DbOption.Pk] = data[pk]
	}
	if data[pk] == nil{ //创建
		return bm.Create(data)
	}else { //更新
		return bm.Update(data, data[pk])
	}
}

//根据PK字段删除记录
func (bm *BaseModel) Delete(id interface{}) (total int64, err error){
	var delIds interface{}
	kind := reflect.TypeOf(id).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		delIds = []interface{}{ id }
	}else{
		delIds = id
	}
	db := bm.DbOption.DB.Table(bm.DbOption.Table).Delete("%s IN ?", bm.DbOption.Pk, delIds)
	return db.RowsAffected, db.Error
}