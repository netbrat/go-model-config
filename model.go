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
type KvsSearchOption struct {
	KvName		string			//kv配置项名
	ExtraWhere	[]interface{}	//额外附加的查询条件
	ReturnPath	bool			//当模型为树型结构时，返回的key是否使用path代替
	NotRowAuth	bool			//是否不使用行级权限过滤条件
	ExtraFields	[]string		//额外附加的查询字段
}

//数据查询选项
type SearchOption struct {
	ExtraWhere 	[]interface{}			//附加的查询条件
	SearchValue	map[string]interface{}	//查询项的值
	ExtraFields	[]string				//额外附加的查询字段
	Page		int						//查询页码
	PageSize	int						//查询记录数
	NotRowAuth	bool					//是否不使用行级权限过滤条件
	NotTotal	bool					//是否不查询总记录数
	NotSearch	bool					//是否不使用配置查询项进行查询
}

//模型结构体
type ConfigModel struct {
	db		*gorm.DB
	config  *Config
}

// 新建一个自定义配制模型
// @param configName string 配制名
func NewConfigModel(c interface{}) (m *ConfigModel, err error) {
	var config *Config

	if reflect.TypeOf(c).Kind() == reflect.Struct {
		config = c.(*Config)
		if err = config.ParseConfig(); err != nil {
			return
		}
	}else {
		//读取模型配置文件
		config, err = GetConfig(c.(string))
		if err != nil {
			return nil, err
		}
	}
	//创建一个连接
	db, err := GetDB(config.ConnName)
	if err != nil {
		return nil, err
	}
	tb := fmt.Sprintf("%s AS %s", config.Table, config.Alias)
	if config.DBName != "" {
		tb = fmt.Sprintf("`%s`.%s", config.DBName, tb)
	}
	m = &ConfigModel{db: db, config: config}
	m.db.Table(tb)
	m.db.Joins(strings.Join(config.Joins," "))
	m.db.Group(strings.Join(m.fieldsAddAlias(config.Groups), ","))
	m.db.Order(strings.Join(m.fieldsAddAlias(config.Orders), ","))

	return
}

// 获取模型配置对象
func (m *ConfigModel) Config() *Config{
	return m.config
}

// 获取数据库连接对象
func (m *ConfigModel) DB() *gorm.DB {
	return m.db
}

// 获取Kv键值列表
func (m *ConfigModel) GetKvs(so *KvsSearchOption) (result map[string]interface{}, err error){
	//检查选项
	if so == nil { so = &KvsSearchOption{KvName: "default"} }
	if so.KvName == "" { so.KvName = "default" }
	if !InArray(so.KvName, m.config.Kvs){
		err =  fmt.Errorf("配置中不存在 [%s] kv 项配置", so.KvName)
		return
	}

	//分析kvs查询的字段
	fields := m.parseKvFields(so.KvName, so.ExtraFields)
	if fields == nil || len(fields) <= 0 {
		return
	}

	//分析kvs查询条件
	db := m.parseWhere(so.ExtraWhere, nil, true, so.NotRowAuth)

	//查询
	var data []map[string]interface{}
	if db.Select(fields).Find(&data); errors.Is(db.Error, gorm.ErrRecordNotFound) {
		err = db.Error
	}

	//处理结果
	result = map[string]interface{}{}
	for _, v := range data {
		key := v["_key"].(string)
		//树形
		if m.config.IsTree  && so.ReturnPath {
			key = v[m.config.TreePathField].(string)
		}
		result[key] = v
	}
	return
}


// 获取数据列表
func (m *ConfigModel) Find(so *SearchOption) (data []map[string]interface{},footer map[string]interface{}, total int64, err error){
	if so == nil {
		so = &SearchOption{}
	}
	//分析查询的字段
	fields, footerFields := m.parseFields(so.ExtraFields)
	if fields == nil || len(fields) <= 0 {
		return
	}
	//分析查询条件
	db := m.parseWhere(so.ExtraWhere, so.SearchValue, so.NotSearch, so.NotRowAuth)

	//分页信息
	offset, limit := GetOffsetLimit(so.Page, so.PageSize)

	//查询
	db.Offset(offset).Limit(limit)
	db.Select(fields).Find(&data)
	if !so.NotTotal { db.Count(&total) }
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		err = db.Error
		return
	}
	//汇总
	if footerFields != nil && len(footerFields) > 0 {
		footer = map[string]interface{}{}
		db.Select(footerFields).Take(&footer)
		if db.Error != nil {
			err = db.Error
			return
		}
	}
	return
}



// 分析查询条件 (此批条件只作用于返回的db对象上，不会作用于模型的db上)
// @param extraWhere 额外的查询条件
// @param searchValues 查询字段值
// @param notSearch 是否使用查询字段条件
// @param notRowAuth 是否使用行级权限进行过滤
func (m *ConfigModel) parseWhere(extraWhere []interface{}, searchValues map[string]interface{}, notSearch bool, notRowAuth bool) *gorm.DB{
	db := m.db.Where("")
	//额外的查询条件
	if extraWhere != nil {
		db.Where(extraWhere[0], extraWhere[1:]...)
	}

	// 模型全局查询条件
	if m.config.Where != "" {
		db.Where(m.config.Where)
	}

	// 模型各查询字段
	if !notSearch{
		if searchValues == nil{
			searchValues = map[string]interface{}{}
		}
		for _, f := range m.config.SearchFields {
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
			db.Where(f.Where, values...)
		}
	}
	if !notRowAuth {

	}
	return db
}

//分析查询字段
// @param	extraFields		额外附加的字段
// @return	fields			最终需要查询的字段名数组
// @return	footerFields	汇总字段
func (m *ConfigModel) parseFields(extraFields []string)(fields []string,footerFields []string){
	fields = make([]string, 0)
	footerFields = make([]string, 0)
	//扩展字段
	fields = append(fields, m.fieldsAddAlias(extraFields)...)
	// 树型必备字段
	if m.config.IsTree {
		treeLevelField := fmt.Sprintf("(LENGTH(%s)/%d) AS __level", m.fieldAddAlias(m.config.TreePathField), m.config.TreePathBit)
		fields = append(fields, treeLevelField)
	}
	for _, f := range m.config.Fields {
		if f.Name == "" || f.Hidden {continue}
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
func (m *ConfigModel) parseKvFields(kvName string, extraFields []string) (fields []string){
	fields = make([]string, 0)
	// kv配置中的字段
	kv, ok := Kv{}, false
	if kv, ok = m.config.Kvs[kvName]; !ok{
		return
	}
	keySep := fmt.Sprintf(",'%s',", kv.KeySep)
	valueSep := fmt.Sprintf(",'%s',", kv.ValueSep)
	keyField := fmt.Sprintf("CONCAT(%s) AS __key", strings.Join(m.fieldsAddAlias(kv.KeyFields), keySep))
	ValueField := fmt.Sprintf("CONCAT(%s) AS __value", strings.Join(m.fieldsAddAlias(kv.ValueFields), valueSep))
	fields = append(fields, keyField, ValueField)

	// 树型必备字段
	if m.config.IsTree {
		treePathField := m.fieldAddAlias(m.config.TreePathField)
		treeLevelField := fmt.Sprintf("(LENGTH(%s)/%d) AS __level", treePathField, m.config.TreePathBit)
		fields = append(fields, treePathField, treeLevelField)
	}
	// 附加字段
	if extraFields != nil {
		fields = append(fields, m.fieldsAddAlias(extraFields)...)
	}
	return
}

// 给字段加表别名
func (m *ConfigModel) fieldAddAlias(field string) string{
	if field == "" { return "" }
	if strings.Contains(field, ".") || strings.Contains(field,"(") {
		return field
	}else{
		return fmt.Sprintf("`%s`.%s", m.config.Alias, strings.Trim(field, " "))
	}
}

// 给字段数组加表别名
func (m *ConfigModel) fieldsAddAlias(fields []string) []string{
	newFields := make([]string, 0)
	for _, v := range fields {
		if v == "" { continue }
		if strings.Contains(v, ".") || strings.Contains(v,"(") {
			newFields = append(newFields, v)
		} else {
			newFields = append(newFields, fmt.Sprintf("`%s`.%s", m.config.Alias,  strings.Trim(v," ")))
		}
	}
	return newFields
}


func (m *ConfigModel) processData(data []map[string]interface{}, footer map[string]interface{}) error{
	if data == nil || len(data) <= 0 { return}
	for _, f := range m.config.Fields {
		if _, ok := data[0][f.Name]; !ok {
			continue
		}
		switch f.Type {
		case FieldTypeEnum: //枚举
			enums := m.config.Enums[f.From]
			for i, _:= range data {
				vString := data[i][f.Name].(string) //字段值
				if f.Multiple{ //多选
					vs := strings.Split(vString, f.Separator)
					newVs := make([]string,0)
					for _, v := range vs{
						newVs = append(newVs, enums[v])
					}
					data[i]["__" + f.Name] = strings.Join(newVs, f.Separator)
				}else{ //单选
					data[i]["__"+ f.Name] = enums[vString]
				}
			}
		case FieldTypeKv: //外联Kv
			joinM, err := NewConfigModel(f.From)
			if err != nil {
				return err
			}
			kvs := joinM.GetKvs(nil)

		}
	}
}