package mc

import (
	"fmt"
	"github.com/netbrat/djson"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

// 模型字段类型
const (
	ModelFieldTypeText 			= "text"		//单行文本（默认）
	ModelFieldTypeAreaText 		= "areatext"	//多行文本
	ModelFieldTypeEnum			= "enum"		//枚举(指定from)
	ModelFieldTypeKv			= "kv"			//外链表Kv对(指定from)
	ModelFieldTypeDate			= "date"		//日期
	ModelFieldTypeDatetime		= "datetime"	//日期时间
)
var fieldTypes []string = []string {ModelFieldTypeText, ModelFieldTypeAreaText, ModelFieldTypeEnum, ModelFieldTypeKv, ModelFieldTypeDate, ModelFieldTypeDatetime}

//字段基础配置对象
type ModelBaseField struct {
	Name        string      `json:"name"`                        //字段名 (必填）
	Title       string      `json:"title"`                       //标题 (默认为name值）
	Type        string      `json:"type" default:"text"`         //字段类型(text:单行文本（默认） | areatext:多行文本 | enum:枚举(指定from) | outjoin:外链表(指定from)|date:日期, datetime:日期时间),如果此字段对应的表格是行权限范围表，则必须使用outjoin
	Info        string      `json:"info"`                        //字段说明（默认""）
	From        string      `json:"from"`                        //指定字段数据来源,当type为enum或outjoin时有效且必填,当为outjoin时填写模型配制ID
	Multiple    bool        `json:"multiple"`                    //是否支持多选（默认false）
	Separator   string      `json:"separator" default:","`       //多选时选项的分隔符，尽量使用默认值逗号，因为mysql的MATCH AGAINST只使用逗号分隔，除非你的业务中不使用些语句（默认逗号,)
	Default     interface{} `json:"default"`                     //默认值（默认""）
	Width       int         `json:"width" default:"120"`         //显示宽度（默认120）
	Height      int         `json:"height" default:"60"`         //多行时的行数，仅当text='areatext'有效 （默认60）
	SelNullText string      `json:"sel_null_text" default:"请选择"` //下拉型默认未选情况下显示的空值文本,当"NO"时不显示，当为""时显示“请选择",默认（"")
	ReturnPath  bool        `json:"return_path"`                 //下拉列表返回路径，而不是KEY，仅针对树型outjoin类型有效 默认 false
}

//字段配置对象
type ModelField struct {
	ModelBaseField
	Editable      *bool  `json:"editable" default:"true"` //是否可编辑（默认true)
	Alias         string `json:"alias"`                   //别名，与SQL中刚好相反，如SQL中：SUM(money) AS total，则此处填写sum(abc)，total为Column单项的key（默认为""）
	Footer        string `json:"footer"`                  //此字段表尾汇总SQL，如SUM(money)，为""，则此字段不汇总
	Filter        *bool  `json:"filter" default:"true"`   //数据保存时是否对内容进行安全过虑(默认true)
	Hidden        bool   `json:"hidden"`                  //是否在列表中隐藏，此列即可单独设置，也会根据权限系统自动进行设置 默认false
	Func          string `json:"func"`                    //值显示的回调
	Align         string `json:"align" default:"left"`    //列表时对齐方式
	Sortable      *bool  `json:"sortable" default:"true"` //列表时是否允许排序（默认true)
	EditHideValue bool   `json:"edit_hide_value"`         //编辑时是否不允许显示原值（默认false)
	Require       bool   `json:"require"`                 //必填字段(默认false)
}

//查询字段配置对象
type ModelSearchField struct {
	ModelBaseField
	Where  string   `json:"where"`  //查询时的条件，表单传过来的值会替换此条件{{this}}字符串，如果是多选，则使用{{inThis}}
	Values []string `json:"values"` //查询时的条件值，默认直接[]
	Br     bool     `json:"br"`     //表单换行显示
}

//键值对配置对象
type ModelKv struct {
	KeyFields   []string `json:"key_fields"`   // 主键（必填）
	ValueFields []string `json:"value_fields"` // 值字段列表 (必填）
	KeySep      string   `json:"value_sep"`    // 多关键字段分隔符（默认_)
	ValueSep    string   `json:"value_sep"`    // 多值字段分隔符（默认_）
	Where       string   `json:"where"`        //查询条件（只作用此kv选择中)
}

//回调js配置对象
type ModelJavascript struct {
	ListStart string `json:"list_start"` //显示列表开始时回调
	ListEnd   string `json:"list_end"`   //显示列表结果时回调
	EditStart string `json:"edit_start"` //编辑弹窗显示回调
	EditEnd   string `json:"edit_end"`   //编辑提交时回调
}

//模型配制属性
type ModelAttr struct {
	Name          string                            `json:"-"`
	ConnName      string                            `json:"conn_name" default:"default"`    //数据库连接名			默认 default
	DBName        string                            `json:"db_name"`                        //数据库名				默认 数据库连接配置中的数据库名
	Table         string                            `json:"table"`                          //数据表名				必填
	Alias         string                            `json:"alias"`                          //表别名 					默认 表名
	Orders        string                          	`json:"orders"`                         //默认排序				选填
	Pk            string                            `json:"pk" default:"id"`                //主键字段名				默认 id
	AutoInc       *bool                             `json:"auto_inc" default:"true"`        //主键是否自增长			默认 true
	UniqueFields  []string                          `json:"unique_fields"`                  //唯一性字段列表			选填
	Where         string                            `json:"where"`                          //基础查询条件 			默认""
	Joins         []string                          `json:"joins"`                          //外联SQL
	Groups        []string                          `json:"groups"`                         //分组SQL
	IsTree        bool                              `json:"is_tree"`                        //是否树型结构表			默认 false
	TreePathBit   int                               `json:"tree_path_bit" default:"2"`      //树型结构路径每层位数		默认 2
	TreePathField string                            `json:"tree_path_field" default:"path"` //树型结构的路径字段		默认 path
	HideCheckbox  bool                              `json:"hide_checkbox" default:"false"`  //列表是否不显示多选框 	默认 false
	Fields        []ModelField                      `json:"fields"`                         //字段列表
	SearchFields  []ModelSearchField                `json:"search_fields"`                  //查询字段列表
	Enums         map[string]map[string]string 		`json:"enums"`                          //枚举列表
	Kvs           map[string]ModelKv                `json:"kvs"`                            //键值对配置结构
	JavaScript    ModelJavascript                   `json:"javascript"`                     //回调js
	//FieldIndexes       map[string]int               `json:"-"`                              //字段索引				填写项
	//SearchFieldIndexes map[string]int               `json:"-"`                              //查询字段索引			非填写项
}

//模型结构体
type Model struct {
	db   *gorm.DB
	attr *ModelAttr
}


// 新建一个自定义配制模型
// @param c  配制名|配制对象
func NewModel(config string) (m *Model) {
	attr := &ModelAttr{}
	if strings.Contains(config, "{") {
		if err := djson.Unmarshal(config, attr, nil); err != nil{
			panic(fmt.Errorf(fmt.Sprintf("解析模型配置出错：%s", err.Error())))
		}
	} else {
		file := fmt.Sprintf("%s%s.json",option.ConfigsFilePath, config)
		if err := djson.FileUnmarshal(file, attr, nil); err != nil {
			panic(fmt.Errorf(fmt.Sprintf("读取模型配置[%s]信息出错：%s", config, err.Error())))
		}
		attr.Name = config
	}
	m = &Model{}
	return m.SetAttr(attr)
}

// 获取模型数据库连接对象本身
func (m *Model) DB() *gorm.DB{
	return m.db
}

// 获取一个新的模型数据库连接对象
func (m *Model) NewDB() *gorm.DB {
	return m.db.Where("")
}

// 获取一个仅包含连接名及表名的连接对象
func (m *Model) BaseDB() *gorm.DB {
	db, err := GetDB(m.attr.ConnName)
	if err != nil {
		panic(err)
	}
	tb := fmt.Sprintf("%s AS %s", m.attr.Table, m.attr.Alias)
	if m.attr.DBName != "" {
		tb = fmt.Sprintf("`%s`.%s", m.attr.DBName, tb)
	}
	db.Table(tb)
	return db
}

func (m *Model) Attr() *ModelAttr {
	return m.attr
}

func (m *Model) SetAttr(attr *ModelAttr) *Model{
	attr.parse()
	m.attr = attr

	//创建一个连接并附加模型基础条件信息
	m.db = m.BaseDB()
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

// 分析配置信息
func (attr *ModelAttr) parse() {
	if attr.ConnName == "" { attr.ConnName = "default"} //如果没有指定数据连接，则使用默认
	if attr.Table == "" { attr.Table = attr.Name } //如果没有指定表名，使用模型配制名称
	if attr.Alias == "" { attr.Alias = attr.Table } //如果表没指定别名，就直接使用表名作别名
	if attr.Pk == "" {attr.Pk = "id"} //没有指定主键，则使用默认id
	if attr.AutoInc == nil { *attr.AutoInc = true} // 没有指定是否自增，使用默认true
	if attr.IsTree{
		if attr.TreePathBit <= 0 { attr.TreePathBit = 2} //没有指定树型路径分隔位数，使用默认2
		if attr.TreePathField == "" { attr.TreePathField = "path"} //没有指定树型路径字段，使用默认path
	}
	// 分析列表字段的基础字段信息
	//attr.FieldIndexes = map[string]int{}
	for i,_ := range attr.Fields {
		f := &attr.Fields[i]
		//attr.FieldIndexes[f.Name] = i
		attr.parseBaseField(&f.ModelBaseField)
		if attr.Fields[i].Editable == nil { *attr.Fields[i].Editable = true}
		if attr.Fields[i].Filter == nil { *attr.Fields[i].Filter = true}
		if attr.Fields[i].Sortable == nil { *attr.Fields[i].Sortable = true}
	}
	// 分析查询字段的基础字段信息
	//attr.SearchFieldIndexes = map[string]int{}
	for i, _ := range attr.SearchFields {
		sf := &attr.SearchFields[i]
		//attr.SearchFieldIndexes[sf.Name] = i
		attr.parseBaseField(&sf.ModelBaseField)
		if sf.Values == nil {
			sf.Values = []string { "?" }
		}
	}
}

// 分析基础字段信息
func (attr *ModelAttr) parseBaseField(field *ModelBaseField){
	// 如果字段没有指定标题，使用字段名
	if field.Title == "" { field.Title = field.Name}
	// 如果指定的字段类型不符，使用默认的text类型
	field.Type = strings.ToLower(field.Type)
	if !inArray(field.Type, fieldTypes) {
		field.Type = ModelFieldTypeText
	}
	// 对from设置值进行分析
	if field.Type == ModelFieldTypeKv {
		if field.From[:1] == ":" {
			field.From = fmt.Sprintf("%s%s",attr.Name, field.From)
		} else if !strings.Contains(field.From,":"){
			field.From = fmt.Sprintf("%s:%s",field.From, "default")
		}
	} else if field.Type == ModelFieldTypeEnum {
		if field.From == "" { field.From = field.Name}
	}
}

func (attr *ModelAttr) getEnum (from string, fType string) (enum map[string]string){

	if fType == ModelFieldTypeKv {
		f := strings.Split(from, ":")
		m := NewModel(f[0])
		kvs, _ := m.FindKvs(&KvsQueryOption{KvName:f[1]})
		enum = map[string]string{}
		for i, v := range kvs{
			enum[i] = cast.ToString(v["__value"])
		}
	}else if fType == ModelFieldTypeEnum {
		enum = attr.Enums[from]
	}
	return
}
