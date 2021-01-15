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

type Enum map[string]map[string]interface{} //Enum 或kvs集

//模型结构体
type Model struct {
	db   *gorm.DB
	attr *ModelAttr
	SearchItems []FormItem
	EditItems	[]FormItem
}


// 新建一个自定义配制模型
// @param configs  配制名
func NewModel(config string) (m *Model) {
	attr := &ModelAttr{}
	if strings.Contains(config, "{") {
		if err := djson.Unmarshal(config, attr, nil); err != nil{
			panic(fmt.Errorf(fmt.Sprintf("解析模型配置出错：%s", err.Error())))
		}
	} else {
		file := fmt.Sprintf("%s%s.json",option.ModelConfigsFilePath, config)
		if err := djson.FileUnmarshal(file, attr, nil); err != nil {
			panic(fmt.Errorf(fmt.Sprintf("读取模型配置[%s]信息出错：%s", config, err.Error())))
		}
		attr.Name = config
	}
	m = &Model{}
	return m.SetAttr(attr)
}


// 获取配置属性
func (m *Model) Attr() *ModelAttr {
	return m.attr
}

// 设置配置属性
func (m *Model) SetAttr(attr *ModelAttr) *Model{
	attr.parse()
	m.attr = attr

	//创建一个连接并附加模型基础条件信息
	m.db = m.BaseDB(true)
	if m.attr.Where != "" {
		m.db.Where(attr.Where)
	}
	if m.attr.Joins != nil || len(m.attr.Joins) > 0 {
		m.db.Joins(strings.Join(attr.Joins, " "))
	}
	if m.attr.Groups != nil || len(m.attr.Groups) > 0 {
		m.db.Group(strings.Join(m.fieldsAddAlias(attr.Groups), ","))
	}
	//m.db.Order(strings.Join(m.fieldsAddAlias(attr.Orders), ","))
	return m
}

// 分析查询项的值，某项不存在，侧使用配置默认值替代
func (m *Model) ParseSearchValues(searchValues map[string]interface{}) (values map[string]interface{}){
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
func (m *Model) GetFromDataMap (from string) (enum Enum){
	enum = make(Enum)
	if from == "" {return}
	isKv := strings.Contains(from, ":")
	if isKv {
		f := strings.Split(from, ":")
		if len(f) < 2 || f[2]=="" {
			f[1] = "default"
		}
		var newM *Model
		if f[0] == m.attr.Name {
			newM = m
		}else{
			newM = NewModel(f[0])
		}
		enum, _ = newM.FindKvs(&KvsQueryOption{KvName:f[1]})
	}else {
		for key, value := range m.attr.Enums[from]{
			enum[key] = map[string]interface{}{
				"__key" : key,
				"__value": value,
			}
		}
	}
	return
}

func (m *Model) ListFields() []*ModelField {
	return m.attr.listFields
}

func (m *Model) FieldIndexMap() map[string]int {
	return m.attr.fieldIndexMap
}
// 查询项
func (m *Model) CreateSearchItems(searchValues map[string]interface{}) {
	values := m.ParseSearchValues(searchValues)
	m.SearchItems = make([]FormItem,0)
	for _, field := range m.attr.SearchFields {
		item := m.createFormItem(&field.ModelBaseField, values[field.Name])
		m.SearchItems = append(m.SearchItems, item)
	}
}

// 编辑项
func  (m *Model) CreateEditItems(values map[string]interface{}) {
	m.EditItems = make([]FormItem,0)
	for i, _ := range m.attr.Fields {
		field := m.attr.Fields[i]
		//如果不允许编辑项（不包含PK字段）
		if !field.Editable {
			continue
		}
		item := m.createFormItem(&field.ModelBaseField, values[field.Name])
		//if field.Name == m.attr.Pk {
		//	value := cast.ToString(values[field.Name])
		//	if m.attr.AutoInc {
		//		item.Html += fmt.Sprintf(`<input type="hidden" id="__%s"  name="__%s" vlaue="%s" />`, field.Name, field.Name, value)
		//	}
		//}
		m.EditItems = append(m.EditItems, item)
	}
}

// 生成单个查询或编辑项
func (m *Model) createFormItem(field *ModelBaseField, value interface{}) FormItem {
	var enum Enum
	if value == nil && field.Default != nil {
		value = field.Default
	}
	// 如果字段是enum或kv，则选读取对应的enum
	if field.From == "" {
		enum = nil
	} else {
		enum = m.GetFromDataMap(field.From)
	}
	item := FormItem{
		Field: field,
		Html:  CreateWidget(field, value, enum),
	}
	return item
}
