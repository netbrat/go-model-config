package mc

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"reflect"
	"strings"
)


//kvs查询选项
type KvsQueryOption struct {
	DB          *gorm.DB      //当此项为空的，使用model.db
	KvName      string        //kv配置项名
	ExtraWhere  []interface{} //额外附加的查询条件
	ReturnPath  bool          //当模型为树型结构时，返回的key是否使用path代替
	NotRowAuth  bool          //是否不使用行级权限过滤条件
	ExtraFields []string      //额外附加的查询字段
}

//数据查询选项
type QueryOption struct {
	DB          *gorm.DB               //当此项为空的，使用model.db
	ExtraWhere  []interface{}          //附加的查询条件
	SearchValue map[string]interface{} //查询项的值
	ExtraFields []string               //额外附加的查询字段
	Order       string               	//排序
	Page        int                    //查询页码
	PageSize    int                    //查询记录数
	NotRowAuth  bool                   //是否不使用行级权限过滤条件
	NotTotal    bool                   //是否不查询总记录数
	NotSearch   bool                   //是否不使用配置查询项进行查询
}



// 获取Kv键值列表
func (m *Model) FindKvs(so *KvsQueryOption) (result map[string]map[string]interface{}, err error){
	//检查选项
	if so == nil { so = &KvsQueryOption{KvName: "default"} }
	if so.KvName == "" { so.KvName = "default" }
	if !inArray(so.KvName, m.attr.Kvs){
		err =  fmt.Errorf("配置中不存在 [%s] kv 项配置", so.KvName)
		return
	}

	//分析kvs查询的字段
	fields := m.parseKvFields(so.KvName, so.ExtraFields)
	if fields == nil || len(fields) <= 0 {
		return
	}

	//分析kvs查询条件
	theDB := m.parseWhere(so.DB, so.ExtraWhere, nil, true, so.NotRowAuth)

	//查询
	var data []map[string]interface{}
	if err = theDB.Select(fields).Find(&data).Error; err !=nil {
		return
	}

	//处理结果
	result = map[string]map[string]interface{}{}
	for _, v := range data {
		key := cast.ToString(v["__key"])
		//树形
		if m.attr.IsTree  && so.ReturnPath {
			key = cast.ToString(v[m.attr.TreePathField])
		}
		result[key] = v
	}
	return
}


// 获取一条数据
func (m *Model) Take(so *QueryOption) (desc map[string]interface{}, exist bool, err error){
	//检查选项
	if so == nil { so = &QueryOption{} }
	//分析查询的字段
	fields, _ := m.parseFields(so.ExtraFields)
	if fields == nil || len(fields) <= 0 { return }
	//分析查询条件
	theDB := m.parseWhere(so.DB, so.ExtraWhere, so.SearchValue, so.NotSearch, so.NotRowAuth)
	desc = make(map[string]interface{})
	err = theDB.Select(fields).Take(&desc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	return desc, true, nil

}

// 获取数据列表
func (m *Model) Find(so *QueryOption) (desc []map[string]interface{},footer map[string]interface{}, total int64, err error){
	//检查选项
	if so == nil { so = &QueryOption{} }
	//分析查询的字段
	fields, footerFields := m.parseFields(so.ExtraFields)
	if fields == nil || len(fields) <= 0 { return }
	//分析查询条件
	theDB := m.parseWhere(so.DB, so.ExtraWhere, so.SearchValue, so.NotSearch, so.NotRowAuth)

	//排序
	if so.Order != ""{
		theDB.Order(so.Order)
	}else if m.attr.Orders != "" {
		theDB.Order(m.attr.Orders)
	}

	//分页信息
	offset, limit := getOffsetLimit(so.Page, so.PageSize)

	//查询
	theDB.Offset(offset).Limit(limit)
	theDB.Select(fields).Find(&desc)
	if !so.NotTotal { theDB.Count(&total) }
	if theDB.Error != nil {
		err = theDB.Error
		return
	}
	//汇总
	if footerFields != nil && len(footerFields) > 0 {
		footer = map[string]interface{}{}
		theDB.Select(footerFields).Take(&footer)
		if theDB.Error != nil {
			err = theDB.Error
			return
		}
	}
	err = m.processData(desc)
	return
}


// 判断是否已有重复数据
func (m *Model) IsUnique(data map[string]interface{})(exist bool, err error){
	//如果没有设置唯一字段，且主键是自增时，直接返回不重复
	if (m.attr.UniqueFields == nil || len(m.attr.UniqueFields) <= 0) && *m.attr.AutoInc {
		return
	}
	db := m.BaseDB()

	for _, field := range m.attr.UniqueFields {
		db.Where(fmt.Sprintf("%s = ?", field), data[field])
	}
	if !*m.attr.AutoInc {
		db.Or(fmt.Sprintf("%s = ?", m.attr.Pk), data[m.attr.Pk])
	}
	var total int64
	if err = db.Count(&total).Error; total > 0{
		exist = true
	}
	return
}

// 更新记录
func (m *Model) Update(data map[string]interface{}, pkValue interface{}) (rowsAffected int64, err error){
	exist := false
	if exist, err = m.IsUnique(data); err != nil{
		return
	}else if exist{
		err = fmt.Errorf("记录已存在")
	}
	db := m.BaseDB()
	db.Where(fmt.Sprintf("%s = ?", m.attr.Pk), pkValue).Updates(data)
	return db.RowsAffected, db.Error
}

// 创建记录
func (m *Model) Create(data map[string]interface{}) (rowsAffected int64, err error) {
	exist := false
	if exist, err = m.IsUnique(data); err != nil{
		return
	}else if exist{
		err = fmt.Errorf("记录已存在")
	}
	db := m.BaseDB().Create(data)
	return db.RowsAffected, db.Error
}

//保存记录（根据pk自动分析是update 或 create）
func (m *Model) Save(data map[string]interface{})(rowsAffected int64, err error)  {
	pk := ""
	if *m.attr.AutoInc { //pk自增表
		pk = m.attr.Pk
	}else{
		pk = "__" + m.attr.Pk
	}
	if data[pk] == nil{ //创建
		return m.Create(data)
	}else { //更新
		return m.Update(data, data[pk])
	}
}

//根据PK字段删除记录
func (m *Model) Delete(id interface{}) (total int64, err error){
	var delIds interface{}
	kind := reflect.TypeOf(id).Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		delIds = id
	}else{
		delIds = []interface{}{ id }
	}
	db := m.BaseDB().Delete(fmt.Sprintf("%s IN ?", m.attr.Pk), delIds)
	return db.RowsAffected, db.Error
}



// 分析查询条件 (此批条件只作用于返回的db对象上，不会作用于模型的db上)
// @param extraWhere 额外的查询条件
// @param searchValues 查询字段值
// @param notSearch 是否使用查询字段条件
// @param notRowAuth 是否使用行级权限进行过滤
func (m *Model) parseWhere(db *gorm.DB, extraWhere []interface{}, searchValues map[string]interface{}, notSearch bool, notRowAuth bool) *gorm.DB{
	var theDB *gorm.DB
	if db == nil {
		theDB = m.NewDB()
	}else {
		theDB = db.Where("")
	}
	//额外的查询条件
	if extraWhere != nil {
		theDB.Where(extraWhere[0], extraWhere[1:]...)
	}


	// 模型各查询字段
	if !notSearch{
		if searchValues == nil{
			searchValues = map[string]interface{}{}
		}
		for _, f := range m.attr.SearchFields {
			// 该查询字段未带条件配置，跳过
			if f.Where == "" {
				continue
			}
			// 未传入查询值时，使用默认值
			if cast.ToString(searchValues[f.Name]) == "" {
				if f.Default == nil {
					continue
					//delete(searchValues, f.Name)
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
			theDB.Where(f.Where, values...)
		}
	}
	if !notRowAuth {

	}
	return theDB
}

//分析查询字段
// @param	extraFields		额外附加的字段
// @return	fields			最终需要查询的字段名数组
// @return	footerFields	汇总字段
func (m *Model) parseFields(extraFields []string)(fields []string,footerFields []string){
	fields = make([]string, 0)
	footerFields = make([]string, 0)
	//扩展字段
	fields = append(fields, m.fieldsAddAlias(extraFields)...)
	// 树型必备字段
	if m.attr.IsTree {
		treeLevelField := fmt.Sprintf("(LENGTH(%s)/%d) AS __level", m.fieldAddAlias(m.attr.TreePathField), m.attr.TreePathBit)
		fields = append(fields, treeLevelField)
	}
	for _, f := range m.attr.Fields {
		if f.Name == "" || f.Hidden { continue } //字段名为空或隐藏字段，跳过
		//基础字段
		field := ""
		if f.Alias == ""{
			field = m.fieldAddAlias(f.Name)
		} else if f.Alias != "" {
			field = fmt.Sprintf("%s AS %s", f.Alias, f.Name)
		}
		fields = append(fields, field)

		//汇总字段
		if f.Footer != "" {
			footerFields = append(footerFields, fmt.Sprintf("%s AS %s", f.Footer, f.Name))
		}
	}
	return
}

// 分析kv字段数组 （仅对通过NewConfigModel创建的模型有效）
// @param 	kvName  		kv配置项名
// @param	extraFields		额外附加的字段
// @return	fields			最终需要查询的KV字段名数组
func (m *Model) parseKvFields(kvName string, extraFields []string) (fields []string){
	fields = make([]string, 0)
	// kv配置中的字段
	kv, ok := ModelKv{}, false
	if kv, ok = m.attr.Kvs[kvName]; !ok{
		return
	}
	keySep := fmt.Sprintf(",'%s',", kv.KeySep)
	valueSep := fmt.Sprintf(",'%s',", kv.ValueSep)
	keyField := fmt.Sprintf("CONCAT(%s) AS __key", strings.Join(m.fieldsAddAlias(kv.KeyFields), keySep))
	ValueField := fmt.Sprintf("CONCAT(%s) AS __value", strings.Join(m.fieldsAddAlias(kv.ValueFields), valueSep))
	fields = append(fields, keyField, ValueField)

	// 树型必备字段
	if m.attr.IsTree {
		treePathField := m.fieldAddAlias(m.attr.TreePathField)
		treeLevelField := fmt.Sprintf("(LENGTH(%s)/%d) AS __level", treePathField, m.attr.TreePathBit)
		fields = append(fields, treePathField, treeLevelField)
	}
	// 附加字段
	if extraFields != nil {
		fields = append(fields, m.fieldsAddAlias(extraFields)...)
	}
	return
}

// 给字段加表别名
func (m *Model) fieldAddAlias(field string) string{
	if field == "" { return "" }
	if strings.Contains(field, ".") || strings.Contains(field,"(") {
		return field
	}else{
		return fmt.Sprintf("`%s`.%s", m.attr.Alias, strings.Trim(field, " "))
	}
}

// 给字段数组加表别名
func (m *Model) fieldsAddAlias(fields []string) []string{
	newFields := make([]string, 0)
	for _, v := range fields {
		if v == "" { continue }
		if strings.Contains(v, ".") || strings.Contains(v,"(") {
			newFields = append(newFields, v)
		} else {
			newFields = append(newFields, fmt.Sprintf("`%s`.%s", m.attr.Alias,  strings.Trim(v," ")))
		}
	}
	return newFields
}



// 对查询的数据进行处理
func (m *Model) processData(data []map[string]interface{})(err error){
	if data == nil || len(data) <= 0 { return }
	for _, f := range m.attr.Fields {
		if _, ok := data[0][f.Name]; !ok {
			continue
		}
		if f.Type == ModelFieldTypeKv || f.Type == ModelFieldTypeEnum {
			enum := m.attr.getEnum(f.From, f.Type)
			for i, _:= range data {
				vString := cast.ToString(data[i][f.Name]) //字段值
				if f.Multiple{ //多选
					vs := strings.Split(vString, f.Separator)
					newVs := make([]string,0)
					for _, v := range vs{
						newVs = append(newVs, enum[v])
					}
					data[i]["__" + f.Name] = strings.Join(newVs, f.Separator)
				}else{ //单选
					data[i]["__"+ f.Name] = enum[vString]
				}
			}
		}
	}
	return
}