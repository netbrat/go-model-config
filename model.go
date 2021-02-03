package mc

import (
	"fmt"
	"github.com/netbrat/djson"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

type FormItem struct {
	Field *ModelBaseField
	Html string
}

type Kvs map[string]map[string]interface{} //Enum 或kvs集

//模型结构体
type Model struct {
	db          *gorm.DB
	attr        *ModelAttr
}


// 新建一个自定义配制模型
// @param configs  配制名
func NewModel(config string) (m *Model) {
	attr := &ModelAttr{}
	if strings.Contains(config, "{") {
		if err := djson.Unmarshal(config, attr, nil); err != nil {
			panic(fmt.Errorf(fmt.Sprintf("解析模型配置出错：%s", err.Error())))
		}
	} else {
		file := fmt.Sprintf("%s%s.json", option.ModelConfigsFilePath, config)
		if err := djson.FileUnmarshal(file, attr, nil); err != nil {
			panic(fmt.Errorf(fmt.Sprintf("读取模型配置[%s]信息出错：%s", config, err.Error())))
		}
		attr.Name = config
	}
	m = &Model{attr: attr}
	return m.SetAttr(m.attr)
}


// 获取配置属性
func (m *Model) Attr() *ModelAttr {
	return m.attr
}


// 设置配置属性
func (m *Model) SetAttr(attr *ModelAttr) *Model {
	attr.parse()
	//m.attr = attr

	//创建一个连接并附加模型基础条件信息
	m.db = m.BaseDB(true)
	if m.attr.Where != "" {
		m.db.Where(attr.Where)
	}
	if m.attr.Joins != nil || len(m.attr.Joins) > 0 {
		m.db.Joins(strings.Join(attr.Joins, " "))
	}
	if m.attr.Groups != nil || len(m.attr.Groups) > 0 {
		m.db.Group(strings.Join(m.FieldsAddAlias(attr.Groups), ","))
	}
	//m.db.Order(strings.Join(m.fieldsAddAlias(attr.Orders), ","))
	return m
}

// 分析查询项的值，某项不存在，侧使用配置默认值替代
func (m *Model) ParseSearchValues(searchValues map[string]interface{}) (values map[string]interface{}) {
	values = make(map[string]interface{})
	//过滤掉空值
	for key, value := range searchValues {
		if cast.ToString(value) != "" {
			values[key] = value
		}
	}
	// 未传入查询值时，使用默认值
	for _, f := range m.attr.SearchFields {
		if _, ok := values[f.Name]; !ok && f.Default != nil {
			values[f.Name] = f.Default
		}
	}
	return
}

// 获取From来源数据
func (m *Model) GetFromKvs (from string) (kvs Kvs){
	kvs = make(Kvs)
	if from == "" {return}
	fromInfo := m.attr.ParseFrom(from)
	if fromInfo.IsKv {
		var newM *Model
		if fromInfo.FromName == m.attr.Name || fromInfo.FromName == "" {
			newM = m
		}else{
			newM = NewModel(fromInfo.FromName)
		}
		kvs, _ = newM.FindKvs(&KvsQueryOption{KvName:fromInfo.kvName})
	}else {
		for key, value := range m.attr.Enums[fromInfo.FromName]{
			kvs[key] = map[string]interface{}{
				"__mc_key" : key,
				"__mc_value": value,
			}
		}
	}
	return
}



// 获取列表字段集
func (m *Model) ListFields() []*ModelField {
	return m.attr.listFields
}

// 获取列表字段索引map
func (m *Model) FieldIndexMap() map[string]int {
	return m.attr.fieldIndexMap
}

// 获取行权限字段信息map
func (m *Model) RowAuthFieldMap() map[string]ModelFieldFromInfo {
	return m.attr.rowAuthFieldMap
}

func (m *Model) IsRowAuth() bool {
	return m.attr.isRowAuth
}

// 创建查询表单项
func (m *Model) CreateSearchItems(searchValues map[string]interface{}) (items []FormItem) {
	values := m.ParseSearchValues(searchValues)
	items = make([]FormItem, 0)
	for i, _ := range m.attr.SearchFields {
		field := m.attr.SearchFields[i]
		item := m.createFormItem(&field.ModelBaseField, values[field.Name])
		items = append(items, item)
	}
	return
}

// 创建编辑表单项
func  (m *Model) CreateEditItems(values map[string]interface{})  (items []FormItem) {
	items = make([]FormItem, 0)
	for i, _ := range m.attr.Fields {
		field := m.attr.Fields[i]
		//如果不允许编辑项（不包含PK字段）
		if !field.Editable {
			continue
		}
		item := m.createFormItem(&field.ModelBaseField, values[field.Name])
		items = append(items, item)
	}
	return
}

// 创建其它表单项
func (m *Model) CreateOtherItems(key string, values map[string]interface{}) (items []FormItem) {
	items = make([]FormItem, 0)
	for i , _ := range m.attr.OtherFields[key] {
		field := m.attr.OtherFields[key][i]
		item := m.createFormItem(&field, values[field.Name])
		items = append(items, item)
	}
	return
}


// 生成单个查询或编辑项
func (m *Model) createFormItem(field *ModelBaseField, value interface{}) FormItem {
	var kvs Kvs
	if value == nil && field.Default != nil {
		value = field.Default
	}
	if field.Multiple {
		value = strings.Split(cast.ToString(value), field.Separator)
	}
	// 如果字段是enum或kv，则选读取对应的enum
	if field.From == "" {
		kvs = nil
	} else {
		kvs = m.GetFromKvs(field.From)
	}
	item := FormItem{
		Field: field,
		Html:  CreateWidget(field, value, kvs),
	}
	return item
}


