package mc

import (
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
	ExtraFields []string      //额外附加的查询字段
	Order       string        //排序
	TreeIndent  *string       //树型模型节点名称前根据层级加前缀字符
}

//数据查询选项
type QueryOption struct {
	DB                  *gorm.DB               //当此项为空的，使用model.db
	ExtraWhere          []interface{}          //附加的查询条件
	Values              map[string]interface{} //查询项的值
	ExtraFields         []string               //额外附加的查询字段
	Order               string                 //排序
	Page                int                    //查询页码（仅对find有效）
	PageSize            int                    //查询记录数 （仅对find有效）
	NotTotal            bool                   //是否不查询总记录数 （仅对find有效）
	NotSearch           bool                   //是否不使用配置查询项进行查询
	NotFoot             bool                   //是否查询汇总项 （仅对find有效）
	TreeIndent          *string                //树型模型节点名称前根据层级加前缀字符
	NotConvertFromValue bool                   //不转换from值， 默认false(转换）
	AttachFromRealValue bool                   //是否附加kv及enum字段原值
	useModelFiledType  string                 // 取list字段还是edit字段列表 (list|edit)
}

type RowData map[string]interface{}

// 获取模型数据库连接对象本身
// 对此修改会影响模型本身的数据库连接
func (m *Model) DB() *gorm.DB {
	return m.db
}

// 获取一个新的模型数据库连接对象
// 对此修改不会影响模型本身的数据库连接
func (m *Model) NewDB() *gorm.DB {
	return m.db.Session(&gorm.Session{}).Where("")
}

// 获取一个仅包含连接名及表名的连接对象
// param isAs 表是否带别名
func (m *Model) BaseDB(isAs bool) *gorm.DB {
	db := GetDB(m.attr.ConnName)
	if isAs {
		tb := fmt.Sprintf("%s AS %s", m.attr.Table, m.attr.Alias)
		if m.attr.DBName != "" {
			tb = fmt.Sprintf("`%s`.%s", m.attr.DBName, tb)
		}
		db.Table(tb)
	} else {
		db.Table(m.attr.Table)
	}
	return db
}

// 获取Kv键值列表
func (m *Model) FindKvs(qo *KvsQueryOption) (desc Kvs, err error) {
	//检查选项
	if qo == nil {
		qo = &KvsQueryOption{KvName: "default"}
	}
	if qo.KvName == "" {
		qo.KvName = "default"
	}
	if !InArray(qo.KvName, m.attr.Kvs) {
		err = fmt.Errorf("配置中不存在 [%s] kv 项配置", qo.KvName)
		return
	}

	//分析kvs查询的字段
	fields := m.ParseKvFields(qo.KvName, qo.ExtraFields)
	if fields == nil || len(fields) <= 0 {
		return
	}

	//分析kvs查询条件
	theDB := m.ParseWhere(qo.DB, qo.ExtraWhere, nil, true)

	//排序
	if qo.Order != "" {
		theDB.Order(qo.Order)
	} else if m.attr.Kvs[qo.KvName].Order != "" {
		theDB.Order(m.attr.Kvs[qo.KvName].Order)
	} else if m.attr.Order != "" {
		theDB.Order(m.attr.Order)
	}

	//查询
	data := make([]map[string]interface{}, 0)
	if err = theDB.Select(fields).Find(&data).Error; err != nil {
		return
	}

	//处理结果
	desc = make(Kvs)
	for i, v := range data {
		key := cast.ToString(v["__mc_key"])
		//树形
		if m.attr.IsTree && qo.ReturnPath {
			key = cast.ToString(v[m.attr.Tree.PathField])
		}
		indent := ""
		if qo.TreeIndent == nil {
			indent = m.attr.Tree.Indent
		} else {
			indent = *qo.TreeIndent
		}
		if m.attr.IsTree && indent != "" { //树形名称字段加前缀
			data[i]["__mc_value"] = nString(indent, cast.ToInt(data[i]["__mc_level"])-1) + cast.ToString(data[i]["__mc_value"])
		}
		desc[key] = v
	}
	return
}

// 获取一条编辑数据
func (m *Model) TakeForEdit(qo *QueryOption) (desc map[string]interface{}, exist bool, err error) {
	indent := ""
	qo.NotConvertFromValue = true
	qo.NotSearch = true
	qo.TreeIndent = &indent
	qo.useModelFiledType = "edit"
	return m.Take(qo)
}

// 获取一条list数据
func (m *Model) Take(qo *QueryOption) (desc map[string]interface{}, exist bool, err error) {
	qo.PageSize = 1
	qo.Page = 1
	qo.NotTotal = true
	qo.NotFoot = true

	if data, _, _, err := m.Find(qo); err != nil  {
		return nil, false, err
	}else if len(data) < 0 {
		return nil, false, nil
	}else{
		return data[0], true, nil
	}
}

// 获取list数据列表
func (m *Model) Find(qo *QueryOption) (desc []map[string]interface{}, foot map[string]interface{}, total int64, err error) {
	//检查选项
	if qo == nil {
		qo = &QueryOption{}
	}
	//分析查询的字段
	fields, footFields := m.ParseFields(qo)
	if fields == nil || len(fields) <= 0 {
		return
	}
	//分析查询条件
	theDB := m.ParseWhere(qo.DB, qo.ExtraWhere, qo.Values, qo.NotSearch)

	//排序
	if qo.Order != "" {
		theDB.Order(qo.Order)
	} else if m.attr.Order != "" {
		theDB.Order(m.attr.Order)
	}

	//分页信息
	offset, limit := getOffsetLimit(qo.Page, qo.PageSize)

	//查询
	desc = make([]map[string]interface{}, 0)
	db := theDB.Session(&gorm.Session{})
	db.Offset(offset).Limit(limit).Select(fields).Find(&desc)
	if !qo.NotTotal {
		db = theDB.Session(&gorm.Session{})
		db.Count(&total)
	}
	if theDB.Error != nil {
		err = theDB.Error
		return
	}
	//汇总
	if !qo.NotFoot && footFields != nil && len(footFields) > 0 {
		foot = make(map[string]interface{})
		if err = theDB.Select(footFields).Offset(0).Limit(1).Take(&foot).Error; err != nil {
			return
		}
	}
	err = m.ProcessData(desc, qo)
	return
}


// 判断是否已有重复数据
func (m *Model) CheckUnique(data map[string]interface{}, oldPkValue interface{})(err error) {
	//如果没有设置唯一字段，且主键是自增时，直接返回不重复
	if (m.attr.UniqueFields == nil || len(m.attr.UniqueFields) <= 0) && m.attr.AutoInc {
		return
	}
	db := m.BaseDB(true)
	pk := m.FieldAddAlias(m.attr.Pk)

	fileTitles := make([]string, 0)

	if oldPkValue != nil {
		db.Where(fmt.Sprintf("%s <> ?", pk), oldPkValue)
		fileTitles = append(fileTitles, m.attr.Fields[m.attr.fieldIndexMap[pk]].Title)
	}

	where := ""
	whereValue := make([]interface{}, 0)
	//检查唯一字段
	for _, field := range m.attr.UniqueFields {
		if where == "" {
			where += fmt.Sprintf(" %s = ?", m.FieldAddAlias(field))
		} else {
			where += fmt.Sprintf(" AND %s = ?", m.FieldAddAlias(field))
		}
		whereValue = append(whereValue, data[field])
		fileTitles = append(fileTitles, m.attr.Fields[m.attr.fieldIndexMap[field]].Title)
	}

	//非自增PK表，检查PK字段
	if !m.attr.AutoInc {
		if where == "" {
			where = fmt.Sprintf("%s = ?", pk)
		} else {
			where = fmt.Sprintf("( %s ) OR ( %s )", where, fmt.Sprintf("%s = ?", pk))
		}
		whereValue = append(whereValue, data[m.attr.Pk])
	}
	db.Where(where, whereValue...)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return err
	} else if total > 0 {
		return &Result{Message:fmt.Sprintf("记录已存在：【%s】存在重复", strings.Join(fileTitles, "、"))}
	}
	return nil
}

// 检查必填字段
func (m *Model) CheckRequiredValues(data map[string]interface{}) (err error) {
	fieldTitles := make([]string, 0)
	//非自增PK表，检查PK字段
	if !m.attr.AutoInc {
		if cast.ToString(data[m.attr.Pk]) == "" {
			fieldTitles = append(fieldTitles, m.attr.Fields[m.attr.fieldIndexMap[m.attr.Pk]].Title)
		}
	}
	//检查配置中的必填字段
	for _, field := range m.attr.Fields {
		if !field.Required {
			continue
		}
		if cast.ToString(data[field.Name]) == "" {
			fieldTitles = append(fieldTitles, field.Title)
		}
	}
	if len(fieldTitles) > 0 {
		return &Result{Message:fmt.Sprintf("【%s】 字段为必填项", strings.Join(fieldTitles, "、"))}
	}
	return
}

// 更新记录
func (m *Model) Updates(data map[string]interface{}, oldPkValue interface{}) (rowsAffected int64, err error) {

	//检查必填项
	if err = m.CheckRequiredValues(data); err != nil {
		return
	}
	//检查重复记录
	if err = m.CheckUnique(data, oldPkValue); err != nil {
		return
	}
	//更新数据
	db := m.BaseDB(false)
	db.Where(fmt.Sprintf("`%s` = ?", m.attr.Pk), oldPkValue).Updates(data)
	return db.RowsAffected, db.Error
}

// 创建记录
func (m *Model) Create(data map[string]interface{}) (rowsAffected int64, err error) {
	//检查必填项
	if err = m.CheckRequiredValues(data); err != nil {
		return
	}
	//检查重复记录
	if err = m.CheckUnique(data, nil); err != nil {
		return
	}
	//创建数据
	db := m.BaseDB(false).Create(data)
	return db.RowsAffected, db.Error
}

//保存记录（根据pk自动分析是update 或 create）
func (m *Model) Save(data map[string]interface{}, oldPkValue interface{})(rowsAffected int64, err error) {
	if oldPkValue == nil { //创建
		return m.Create(data)
	} else { //更新
		return m.Updates(data, oldPkValue)
	}
}

//根据PK字段删除记录
func (m *Model) Delete(id interface{}) (rowsAffected int64, err error) {
	var delIds interface{}
	kind := reflect.TypeOf(id).Kind()
	symbol := ""
	if kind == reflect.Array || kind == reflect.Slice {
		symbol = "IN"
		delIds = id
	} else {
		symbol = "="
		delIds = []interface{}{id}
	}
	db := m.BaseDB(false).Where(fmt.Sprintf("`%s` %s ?", m.attr.Pk, symbol), delIds).Delete(nil)
	return db.RowsAffected, db.Error
}

// 分析查询条件 (此批条件只作用于返回的db对象上，不会作用于模型的db上)
// @param extraWhere 额外的查询条件
// @param searchValues 查询字段值
// @param notSearch 是否使用查询字段条件
func (m *Model) ParseWhere(db *gorm.DB, extraWhere []interface{}, searchValues map[string]interface{}, notSearch bool) *gorm.DB {
	var theDB *gorm.DB
	if db == nil {
		theDB = m.NewDB()
	} else {
		theDB = db.Where("")
	}
	//额外的查询条件
	if extraWhere != nil {
		theDB.Where(extraWhere[0], extraWhere[1:]...)
	}

	// 模型各查询字段
	if !notSearch {
		searchValues = m.ParseSearchValues(searchValues)
		for _, f := range m.attr.SearchFields {
			// 该查询字段未带条件配置 或 未传值，跳过
			if _, ok := searchValues[f.Name]; !ok {
				continue
			}
			if f.Where == "" {
				f.Where = fmt.Sprintf("%s = ?", m.FieldAddAlias(f.Name))
				f.Values = []string{"?"}
			}
			// 查询值与查询条件匹配
			values := make([]interface{}, 0)
			if f.Between { //范围值
				vType := reflect.TypeOf(searchValues[f.Name]).Kind()
				var vs []string
				if vType == reflect.Array || vType == reflect.Slice {
					vs = searchValues[f.Name].([]string)
				} else {
					vs = strings.Split(cast.ToString(searchValues[f.Name]), f.BetweenSep)
				}
				for i, v := range f.Values {
					if v == "?" {
						values = append(values, vs[i])
					} else {
						values = append(values, strings.ReplaceAll(v, "?", vs[i]))
					}
				}
			} else { //单个值
				for _, v := range f.Values {
					if v == "?" {
						values = append(values, searchValues[f.Name])
					} else {
						values = append(values, strings.ReplaceAll(v, "?", cast.ToString(searchValues[f.Name])))
					}
				}
			}
			theDB.Where(f.Where, values...)
		}
	}
	//受行权限控制的字段进行数据权限过滤
	for fieldName, fromInfo := range m.attr.rowAuthFieldMap {
		if rowAuth, isAllAuth := option.ModelAuth.GetRowAuthCallback(fromInfo.FromName); !isAllAuth {
			theDB.Where(fmt.Sprintf("%s IN ?", m.FieldAddAlias(fieldName)), rowAuth)
		}
	}
	//如果自身也是行权限模型，则进行本身数据权限过滤
	if m.attr.isRowAuth {
		if rowAuth, isAllAuth := option.ModelAuth.GetRowAuthCallback(m.attr.Name); !isAllAuth {
			theDB.Where(fmt.Sprintf("%s IN ?", m.FieldAddAlias(m.attr.Pk)), rowAuth)
		}
	}

	return theDB
}

//分析查询字段
// @param	extraFields		额外附加的字段
// @return	fields			最终需要查询的字段名数组
// @return	footFields	汇总字段
func (m *Model) ParseFields(qo *QueryOption)(fields []string,footFields []string) {
	fields = make([]string, 0)
	footFields = make([]string, 0)
	//扩展字段
	fields = append(fields, m.FieldsAddAlias(qo.ExtraFields)...)
	// 树型必备字段
	if m.attr.IsTree {
		fields = append(fields, m.ParseTreeExtraField()...)
	}
	var modelFields []*ModelField
	if strings.ToLower(qo.useModelFiledType) == "edit" {
		modelFields = m.attr.editFields
	}else{
		modelFields = m.attr.listFields
	}
	for _, field := range modelFields {
		//基础字段
		fieldName := ""
		if field.Alias == "" {
			fieldName = m.FieldAddAlias(field.Name)
		} else if field.Alias != "" {
			fieldName = fmt.Sprintf("%s AS %s", field.Alias, field.Name)
		}
		fields = append(fields, fieldName)

		//汇总字段
		if field.Foot != "" {
			footFields = append(footFields, fmt.Sprintf("%s AS %s", field.Foot, field.Name))
		}
	}

	return
}

// 分析kv字段数组 （仅对通过NewConfigModel创建的模型有效）
// @param 	kvName  		kv配置项名
// @param	extraFields		额外附加的字段
// @return	fields			最终需要查询的KV字段名数组
func (m *Model) ParseKvFields(kvName string, extraFields []string) (fields []string) {
	fields = make([]string, 0)
	// kv配置中的字段
	kv, ok := ModelKv{}, false
	if kv, ok = m.attr.Kvs[kvName]; !ok {
		return
	}
	//keySep := fmt.Sprintf(",'%s',", kv.KeySep)
	//valueSep := fmt.Sprintf(",'%s',", kv.ValueSep)
	keyField := fmt.Sprintf("%s AS `__mc_key`", m.FieldAddAlias(kv.KeyField))
	valueField := fmt.Sprintf("%s AS `__mc_value`", m.FieldAddAlias(kv.ValueField))
	fields = append(fields, keyField, valueField)

	// 树型必备字段
	if m.attr.IsTree {
		treePathField := m.FieldAddAlias(m.attr.Tree.PathField)
		fields = append(append(fields, treePathField), m.ParseTreeExtraField()...)
	}
	// 附加字段
	if extraFields != nil {
		fields = append(fields, m.FieldsAddAlias(extraFields)...)
	}
	return
}

// 给字段加表别名
func (m *Model) FieldAddAlias(field string) string {
	if field == "" {
		return ""
	}
	if strings.Contains(field, ".") || strings.Contains(field, "(") {
		return field
	} else {
		return fmt.Sprintf("`%s`.`%s`", m.attr.Alias, strings.Trim(field, " "))
	}
}

// 给字段数组加表别名
func (m *Model) FieldsAddAlias(fields []string) []string {
	newFields := make([]string, 0)
	for _, v := range fields {
		if v == "" {
			continue
		}
		if strings.Contains(v, ".") || strings.Contains(v, "(") {
			newFields = append(newFields, v)
		} else {
			newFields = append(newFields, fmt.Sprintf("`%s`.`%s`", m.attr.Alias, strings.Trim(v, " ")))
		}
	}
	return newFields
}


// 对查询的数据进行处理
func (m *Model) ProcessData(data []map[string]interface{}, qo *QueryOption)(err error) {
	if data == nil || len(data) <= 0 {
		return
	}

	//序号
	if m.attr.Number {
		for i, _ := range data {
			data[i]["__mc_index"] = (qo.Page -1) * qo.PageSize + i + 1
		}
	}
	//转换成from值
	if !qo.NotConvertFromValue {
		for _, f := range m.attr.Fields {
			if _, ok := data[0][f.Name]; !ok {
				continue
			}
			if f.From != "" {
				enum := m.GetFromKvs(f.From)
				for i, _ := range data {
					if qo.AttachFromRealValue { //附加字段原值真实值
						data[i]["__mc_"+f.Name] = data[i][f.Name]
					}
					vString := cast.ToString(data[i][f.Name]) //字段值
					if f.Multiple {                           //多选
						vs := strings.Split(vString, f.Separator)
						newVs := make([]string, 0)
						for _, v := range vs {
							newVs = append(newVs, cast.ToString(enum[v]["__mc_value"]))
						}
						data[i][f.Name] = strings.Join(newVs, f.Separator)
					} else { //单选
						data[i][f.Name] = cast.ToString(enum[vString]["__mc_value"])
					}
				}
			}
		}
	}
	//树形
	indent := ""
	if qo.TreeIndent == nil {
		indent = m.attr.Tree.Indent
	} else {
		indent = *qo.TreeIndent
	}
	if m.attr.IsTree && indent != "" { //树形名称字段加前缀
		for i, _ := range data {
			data[i][m.attr.Tree.NameField] = nString(indent, cast.ToInt(data[i]["__mc_level"])-1) + cast.ToString(data[i][m.attr.Tree.NameField])
		}
	}
	return
}


// 分析树形结构查询必须的扩展字段
func (m *Model) ParseTreeExtraField() (field []string) {
	pathField := m.FieldAddAlias(m.attr.Tree.PathField)
	__mc_pathField := fmt.Sprintf("`__mc_%s`.`%s`", m.attr.Table, m.attr.Tree.PathField)
	__mc_pkField := fmt.Sprintf("`__mc_%s`.`%s`", m.attr.Table, m.attr.Pk)

	field = make([]string, 3)
	//层级字段
	field[0] = fmt.Sprintf("CEILING(LENGTH(%s)/%d) AS `__mc_level`", pathField, m.attr.Tree.PathBit)
	//父节点字段
	field[1] = fmt.Sprintf("(SELECT %s FROM `%s` AS `__mc_%s` WHERE %s=LEFT(%s, LENGTH(%s)-%d) LIMIT 1) AS `__mc_parent`",
		__mc_pkField, m.attr.Table, m.attr.Table, __mc_pathField, pathField, pathField, m.attr.Tree.PathBit)
	//字节点数字段
	field[2] = fmt.Sprintf("(SELECT count(%s) FROM `%s` AS `__mc_%s` WHERE %s=LEFT(%s, LENGTH(%s)-%d) LIMIT 1) AS `__mc_child_count`",
		__mc_pkField, m.attr.Table, m.attr.Table, pathField, __mc_pathField, __mc_pathField, m.attr.Tree.PathBit)
	return
}