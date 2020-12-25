package mc

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"strings"
)

type KvsSearchOption struct {
	DB			*gorm.DB		//使用的数据查询对象，如为nil时，使用model内置db对象
	KvName		string			//kv配置项名
	ExtraWhere	interface{}		//附加的查询条件
	returnPath	bool			//当模型为树型结构时，返回的key是否使用path代替
	NotRowAuth	bool			//是否不使用行级权限过滤条件
	Indent		string			//当模型为树型结构时，层级缩进符, 空字符串时不缩进
	ExtraFields	[]string		//额外附加的查询字段
}

//模型结构体
type model struct {
	db 				*gorm.DB
	config 			*Config
	configName		string
	connName 		string
	dbName			string
	table			string
	pkField			string
	autoIncrement 	bool
	uniqueFields	[]string
}

// 新建一个模型
// @param connName 		数据库连接配置名
// @param table 		数据表名
// @param pkField		主键字段名	 	默认 id
// @param autoIncrement	主键是否自增长
// @param uniqueFields	唯一性字段列表 （不包括主键，可为nil)
func NewModel(connName string, table string, pkField string, autoIncrement bool, uniqueFields []string) *model {
	if connName == "" || table == "" {
		panic("connName 和 table 不能为空字符串")
	}
	if pkField == "" {
		pkField = "id"
	}
	m := &model{
		db: 			GetDB(connName),
		connName: 		connName,
		dbName:			GetDbNameByConfig(connName),
		table: 			table,
		pkField: 		pkField,
		autoIncrement: 	autoIncrement,
		uniqueFields: 	uniqueFields,
	}
	return m
}

// 新建一个自定义配制模型
// @param configName string 配制名
func NewConfigModel(configName string) *model {
	//读取模型配置文件
	config, err := GetConfig(configName)
	if err != nil {
		panic(err)
	}
	m := NewModel(config.ConnName, config.Table, config.PkField, config.AutoIncrement, config.UniqueFields)
	m.config = config
	m.connName = configName
	if config.DbName != "" {
		m.dbName = config.DbName
	}


	return m
}

// 获取自自定义配制信息
func (m *model) Config() Config {
	return *m.config
}

// 获取自定义配制名
func (m *model) ConfigName() string {
	return m.configName
}

// 获取主键字段名
func (m *model) PkField() string{
	return m.pkField
}

// 获取主键是否自增长
func (m *model) AutoIncrement() bool {
	return m.autoIncrement
}

// 获取唯一字段数组
func (m *model) UniqueFields() []string {
	return m.uniqueFields
}

// 获取数据库名
func (m *model) DbName() string {
	return m.dbName
}

// 获取数据库连接配制名
func (m *model) ConnName() string {
	return m.connName
}

// 获取数据库连接对象
func (m *model) DB() *gorm.DB {
	return m.db
}

// 设置查询条件
func (m *model) Where(query interface{}, args ...interface{}) *model{
	m.db = m.db.Where(query, args...)
	return m
}

// 设置外联
func (m *model) Joins(query string, args ...interface{}) *model {
	m.db = m.db.Joins(query, args...)
	return m
}

// 设置分组
func (m *model) Group(query string) *model {
	m.db = m.db.Group(query)
	return m
}

// 设置having
func (m *model) Having(query interface{}, values ...interface{}) *model {
	m.db = m.db.Having(query, values...)
	return m
}

// 设置排序
func (m *model) Order(value interface{}, reorder ...bool) *model {
	m.db = m.db.Order(value, reorder...)
	return m
}


// 获取Kv键值列表
func (m *model) GetKvs(so *KvsSearchOption) (data []map[string]interface{}, err error){
	if so.KvName == "" { so.KvName = "default"}
	if !InArray(so.KvName,m.config.Kvs){
		return nil, fmt.Errorf("配置中不存在 [%s] kv 项配置", so.KvName)
	}
	db := m.db
	if so.DB != nil {
		db = so.DB
	}
	db = m.parseConfigWhere(so.ExtraWhere, nil, true, false)
	fields := m.parseConfigKvFields(so.KvName, so.ExtraFields)
	if db = db.Select(fields).Find(&data); !db.RecordNotFound() {
		err = db.Error
	}
	return
}


// 分析查询条件（仅对通过NewConfigModel创建的模型有效）
// @param extraWhere 额外的查询条件
// @param searchValues 查询字段值
// @param notSearch 是否使用查询字段条件
// @param notRowAuth 是否使用行级权限进行过滤
func (m *model) parseConfigWhere(extraWhere interface{}, searchValues map[string]interface{}, notSearch bool, notRowAuth bool) (rdb *gorm.DB){
	rdb = m.db
	// 模型全局查询条件
	if m.config.BaseSearch.Where != "" {
		rdb = rdb.Where(m.config.BaseSearch.Where)
	}
	//额外的查询条件
	if extraWhere != "" {
		rdb = rdb.Where(extraWhere)
	}
	// 模型各查询字段
	if !notSearch{
		for _, f := range m.config.SearchFields {
			// 该查询字段未带条件配置，跳过
			if f.Where == "" {
				continue
			}
			// 未传入查询值时，使用默认值
			if cast.ToString(searchValues[f.Name]) == "" {
				if f.Default != nil {
					delete(searchValues, f.Name)
				} else {
					searchValues[f.Name] = f.Default
				}
			}
			// 查询值与查询条件匹配
			values := make([]interface{}, 0)
			for _, v := range f.Values {
				if v == "?" {
					values = append(values, searchValues[f.Name])
				} else {
					values = append(values, strings.ReplaceAll(v, "?", cast.ToString(searchValues[f.Name])))
				}
			}
			rdb = rdb.Where(f.Where, values...)
		}
	}
	if !notRowAuth {

	}
	return
}


// 分析kv字段数组 （仅对通过NewConfigModel创建的模型有效）
// @param 	kvName  kv配置项名
// @return	fields	[]string		最终需要查询的KV字段名数组
func (m *model) parseConfigKvFields(kvName string, extraFields []string) (fields []string){
	fields = make([]string, 0)

	// kv配置中的字段
	kv, ok := Kv{}, false
	if kv, ok = m.config.Kvs[kvName]; !ok{
		return
	}
	keySep := fmt.Sprintf("'%s',%s.", kv.KeySep, m.config.BaseSearch.Alias)
	valueSep := fmt.Sprintf("'%s',%s.", kv.ValueSep, m.config.BaseSearch.Alias)
	keyField := fmt.Sprintf("CONCAT(%s.%s) AS _key", m.config.BaseSearch.Alias, strings.Join(kv.KeyFields, keySep))
	ValueField := fmt.Sprintf("CONCAT(%s.%s) AS _value", m.config.BaseSearch.Alias, strings.Join(kv.ValueFields, valueSep))
	fields = append(fields, keyField, ValueField)

	// 树型必备字段
	if m.config.IsTree {
		fields = append(fields, m.config.TreePathField, m.config.TreeLevelField)
	}
	// 附加字段
	if extraFields != nil {
		fields = append(fields, extraFields...)
	}
	return
}








//


//type SearchOption struct {
//	Where			string					//查询条件
//	WhereValue		[]interface{}			//查询值
//	Fields			[]string				//查询字段
//	Page			int						//查询页码
//	PageSize 		int						//分页大小
//	OrderBy			string					//排序
//	Join			string					//外联
//	Group			string					//分组
//	Alias			string					//别名
//	Having			string
//	NotTotal		bool					//是否不查询总记录数
//	NotRowAuth		bool					//是否不使用行级权限,默认为true
//	SearchValues 	map[string]interface{}	//查询字段值
//	IsSearch		bool					//是否使用查询字段进行查询
//}
//
//func (so SearchOption) Offset() int{
//	return (so.Page-1) * so.PageSize
//}
//
//type BaseModel struct {
//	DbOption 	*DbOption
//}
//
////设置数据库操作选项
//func (bm *BaseModel) SetDbOption(connName string, dbName string, table string, pk string, autoIncrement bool, uniqueFields []string){
//	bm.DbOption = &DbOption{}
//	bm.DbOption.Set(connName, dbName, table, pk, autoIncrement, uniqueFields)
//}
//
////获取单条记录
//func (bm *BaseModel) First(so SearchOption)(data map[string]interface{}, err error){
//	err = bm.DbOption.DB.
//		Table(bm.DbOption.Table).
//		Order(so.OrderBy).
//		Select(so.Fields).
//		Where(so.Where,so.WhereValue...).
//		Joins(so.Join).
//		Having(so.Having).
//		First(data).Error
//	return
//}
//
////获取记录集
//func (bm *BaseModel) Find(so SearchOption)(data []map[string]interface{}, total int, err error){
//	db := bm.DbOption.DB.
//		Table(bm.DbOption.Table).
//		Order(so.OrderBy).
//		Select(so.Fields).
//		Where(so.Where,so.WhereValue...).
//		Joins(so.Join).
//		Limit(so.PageSize).
//		Offset(so.Offset()).
//		Having(so.Having).
//		Find(data)
//	if err = db.Error; err != nil{
//		return
//	}
//	if !so.NotTotal {
//		err = db.Count(&total).Error
//	}
//	return
//}
//
////判断记录是否存在
//func (bm *BaseModel) IsExist(data map[string]interface{}) (exist bool, err error){
//	where := ""
//	whereValue := make([]interface{},0)
//	db := bm.DbOption.DB.Table(bm.DbOption.Table)
//	for _,v := range bm.DbOption.UniqueFields {
//		if where == "" {
//			where += fmt.Sprintf(" AND %s = ?", v)
//		}else{
//			where = fmt.Sprintf("%s = ?", v)
//		}
//		whereValue = append(whereValue, data[v])
//	}
//
//	if !bm.DbOption.AutoIncrement{
//		where = fmt.Sprintf("(%s) OR (%s = ?)", where, bm.DbOption.Pk)
//		whereValue = append(whereValue, data[bm.DbOption.Pk])
//	}
//	total := 0
//	db = db.Where(where, whereValue...).Count(&total)
//	if total >0 {
//		exist = true
//	}
//	return exist, db.Error
//}
//
////更新记录
//func (bm *BaseModel) Update(data map[string]interface{}, id interface{})(total int64, err error){
//	exist := false
//	if exist, err = bm.IsExist(data); err != nil{
//		return
//	}else if exist {
//		err = fmt.Errorf("记录已存在")
//		return
//	}
//	where := fmt.Sprintf("%s = ?", bm.DbOption.Pk)
//	db := bm.DbOption.DB.Table(bm.DbOption.Table).Where(where, id).Update(data)
//	return db.RowsAffected, db.Error
//}
//
////创建记录
//func (bm *BaseModel) Create(data map[string]interface{})(total int64, err error){
//	exist := false
//	if exist, err = bm.IsExist(data); err != nil{
//		return
//	}else if exist {
//		err = fmt.Errorf("记录已存在")
//		return
//	}
//	db := bm.DbOption.DB.Table(bm.DbOption.Table).Create(data)
//	return db.RowsAffected, db.Error
//}
//
////保存记录（根据pk自动分析是update 或 create）
//func (bm *BaseModel) Save(data map[string]interface{})(total int64, err error){
//	pk := ""
//	where := map[string]interface{}{}
//	if bm.DbOption.AutoIncrement { //pk自增表
//		pk = bm.DbOption.Pk
//	}else{
//		pk = "__" + bm.DbOption.Pk
//		where[bm.DbOption.Pk] = data[pk]
//	}
//	if data[pk] == nil{ //创建
//		return bm.Create(data)
//	}else { //更新
//		return bm.Update(data, data[pk])
//	}
//}
//
////根据PK字段删除记录
//func (bm *BaseModel) Delete(id interface{}) (total int64, err error){
//	var delIds interface{}
//	kind := reflect.TypeOf(id).Kind()
//	if kind != reflect.Array && kind != reflect.Slice {
//		delIds = []interface{}{ id }
//	}else{
//		delIds = id
//	}
//	db := bm.DbOption.DB.Table(bm.DbOption.Table).Delete("%s IN ?", bm.DbOption.Pk, delIds)
//	return db.RowsAffected, db.Error
//}
//
//// @title parseWhere
//// @description 分析查询条件
//// @param	db		*gorm.DB
//// @param 	so 		SearchOption	模型查询字段各项的值
//// @return	rdb		*gorm.DB		附带最终的查询条件的db对象
//func (bm *BaseModel) parseWhere(db *gorm.DB, so SearchOption) (rdb *gorm.DB) {
//	rdb = db
//	if so.Where != "" {
//		rdb = rdb.Where(so.Where, so.WhereValue...)
//	}
//	return
//}