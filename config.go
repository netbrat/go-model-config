package mc

import (
	"bytes"
	"fmt"
	"github.com/netbrat/go-model-config/mcjson"
	"io/ioutil"
)

/**
自定义模型配置结构定义

说明：如果一个表的某些字段是外联到一些行级权限表，则不论上字段是否显示都必须配置到fields上，且类型必须是outjoin

*/

const (
	FieldTypeText 			= "text"		//单行文本（默认）
	FieldTypeAreaText 		= "areatext"	//多行文本
	FieldTypeEnum			= "enum"		//枚举(指定from)
	FieldTypeOutJoin		= "outjoin"		//外链表(指定from)
	FieldTypeDate			= "date"		//日期
	FieldTypeDatetime		= "datetime"	//日期时间
)

//字段基础配置结构
type BaseField struct{
	Name			string	`json:"name"`			//字段名 (必填）
	Title 			string 	`json:"title"`			//标题 (默认为name值）
	Type 			string	`json:"type"`			//字段类型(text:单行文本（默认） | areatext:多行文本 | enum:枚举(指定from) | outjoin:外链表(指定from)|date:日期, datetime:日期时间),如果此字段对应的表格是行权限范围表，则必须使用outjoin
	Info 			string	`json:"info"`			//字段说明（默认""）
	From 			string	`json:"from"`			//指定字段数据来源,当type为enum或outjoin时有效且必填,当为outjoin时填写模型配制ID
	Multiple 		bool	`json:"multiple"`		//是否支持多选（默认false）
	Separator 		string	`json:"separator"`		//多选时选项的分隔符，尽量使用默认值逗号，因为mysql的MATCH AGAINST只使用逗号分隔，除非你的业务中不使用些语句（默认逗号,)
	DefValue		string	`json:"def_value"`		//默认值（默认""）
	Width			int		`json:"width"`			//显示宽度（默认120）
	Height			int		`json:"height"`			//多行时的行数，仅当text='multitext'有效 （默认60）
	SelNullText 	string	`json:"sel_null_text"`	//下拉型默认未选情况下显示的空值文本,当"NO"时不显示，当为""时显示“请选择",默认（"")
	ReturnPath 		bool	`json:"return_path"` 	//下拉列表返回路径，而不是KEY，仅针对树型outjoin类型有效
}

//字段配置结构
type Field struct {
	BaseField
	Editable		bool  `json:"editable"`			//是否可编辑（默认true)
	Alias			string	`json:"alias"`				//别名，与SQL中刚好相反，如SQL中：SUM(money) AS total，则此处填写sum(abc)，total为Column单项的key（默认为""）
	Footer			string	`json:"footer"`				//此字段表尾汇总SQL，如SUM(money)，为""，则此字段不汇总
	NoFilter		bool	`json:"no_filter"`			//数据保存时不对内容进行安全过虑(默认false)
	Hidden			bool	`json:"hidden"`				//是否在列表中不显示，此列即可单独设置，也会根据权限系统自动进行设置（默认false)
	Func			string	`json:"func"`				//值显示的回调
	Align			string	`json:"align"`				//列表时对齐方式
	NoSortable		bool	`json:"sortable"`			//列表时是否允许排序（默认false)
	EditHideValue 	bool	`json:"edit_hide_value"` 	//编辑时是否不允许显示原值（默认false)
	Require			bool	`json:"require"`			//必填字段(默认false)
}

//查询字段配置结构
type SearchField struct {
	BaseField
	Where		string		`json:"where"`			//查询时的条件，表单传过来的值会替换此条件{{this}}字符串，如果是多选，则使用{{inThis}}
	Values		[]string	`json:"values"`			//查询时的条件值，默认直接[]
	Br			string		`json:"br"`				//表单换行显示
}


//基础查询配置结构
type BaseSearch struct {
	Where	string			`json:"where"`			//基础查询条件 (默认"")
	Alias	string			`json:"alias"`			//表别名 （默认表名）
	Join	string			`json:"join"`			//外联SQL
	Group	string			`json:"group"`			//分组SQL
}

//键值对配置结构
type Kv struct {
	KeyFields		[]string		`json:"key_fields"`				// 主键（必填）
	ValueFields		[]string		`json:"value_fields"`			// 值字段列表 (必填）
	KeyConnector	string			`json:"value_connector"`		//多关键字段连接符（默认_)
	ValueConnector 	string			`json:"value_connector"`		// 多值字段连接符（默认_）
}

//回调js配置结构
type JavaScript struct {
	ListStart	string		`json:"list_start"`			//显示列表开始时回调
	ListEnd		string		`json:"list_end"`			//显示列表结果时回调
	EditStart	string		`json:"edit_start"`			//编辑弹窗显示回调
	EditEnd		string		`json:"edit_end"`			//编辑提交时回调
}



//自定义模型整体配置结构
type Config struct {
	Name				string						`json:"-"`
	ConnName			string						`json:"conn_name" default:"aaa"`			//数据库连接名(默认：default)
	DbName				string						`json:"db_name"`			//数据库名(默认：数据库连接配置中的数据库名)
	Table				string						`json:"table"`				//数据表名
	Pk 					*string						`json:"pk" default:"ccc"`					//主键Id
	AutoIncrement 		*bool						`json:"auto_increment"`		//主键是否自增长（默认true)
	OrderBy				string						`json:"order_by"`			//默认排序
	IsTree				bool						`json:"is_tree"`			//是否树型结构表
	TreePathBit			int							`json:"tree_path_bit"`		//树型结构路径每层位数
	TreeLevelField		string						`json:"tree_level_field"`	//树型结构的层级字段
	TreePathField		string						`json:"tree_path_field"`	//树型结构的路径字段
	ShowCheck			bool						`json:"show_check"`			//列表是否显示多选框 (默认 true)
	FieldIndexes		map[string]int				`json:"-"`					//字段索引
	Fields				[]Field						`json:"fields"`				//字段列表
	SearchFieldIndexes	map[string]int				`json:"-"`					//查询字段索引
	SearchFields 		[]SearchField				`json:"search_fields"`		//查询字段列表
	Enums				map[string]interface{}		`json:"enums"`				//枚举列表
	BaseSearch			*BaseSearch					`json:"base_search"`		//基础查询信息
	Kvs					map[string]Kv				`json:"kvs"`				//键值对配置结构
	JavaScript			JavaScript					`json:"javascript"`			//回调js
}


func GetConfig(mcName string) (mc *Config, err error) {
	file := fmt.Sprintf("%s%s.json",option.ConfigsFilePath, mcName)
	mc = &Config{}
	data := []byte{}
	if data, err = ioutil.ReadFile(file); err != nil{
		return
	}else{
		data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
		if err = mcjson.JsonToStruct(data, mc); err != nil {
			return
		}
	}

	mc.Name = mcName
	fmt.Println(mc.BaseSearch)
	//mc.parseConfig()
	return
}

//func (mc *Config) parseConfig() {
//	if mc.ConnName == "" { mc.ConnName = "default" }
//	if mc.DbName == "" { mc.DbName = "" }
//	if mc.Table == "" { mc.Table = mc.Name }
//	if mc.Pk == "" { mc.Pk = "id" }
//	if mc.AutoIncrement == nil { *mc.AutoIncrement = true }
//	if mc.IsTree {
//		if mc.TreePathBit <= 0 { mc.TreePathBit = 2 }
//		if mc.TreeLevelField == "" { mc.TreeLevelField = "level" }
//		if mc.TreePathField == "" { mc.TreePathField = "path" }
//	}
//	if mc.ShowCheck == nil { *mc.ShowCheck = true}
//}


//func (mc *Config) parseFields() {
//	for k, _ := range mc.Fields{
//		if mc.Fields[k].Title == "" {mc.Fields[k].Title = mc.Fields[k].Name}
//		if mc.Fields[k].Type == "" {mc.Fields[k].Title = FieldTypeText}
//		if mc.Fields[k].Multiple && mc.Fields[k].Separator == "" { mc.Fields[k].Separator = ","}
//		if mc.Fields[k].SelNullText == "" {mc.Fields[k].SelNullText = "请选择"}
//		if mc.Fields[k].Width <=0 { mc.Fields[k].Width = 120}
//		if mc.Fields[k].Height <=0 { mc.Fields[k].Height = 60}
//		if mc.Fields[k].Editable == nil {*mc.Fields[k].Editable = true}
//	}
//}
//
//func (mc *Config) parseSearchFields() {
//	for k, _ := range mc.SearchFields{
//		if mc.SearchFields[k].Title == "" {mc.SearchFields[k].Title = mc.Fields[k].Name}
//		if mc.SearchFields[k].Type == "" {mc.SearchFields[k].Title = FieldTypeText}
//		if mc.SearchFields[k].Multiple && mc.SearchFields[k].Separator == "" { mc.Fields[k].Separator = ","}
//		if mc.SearchFields[k].SelNullText == "" {mc.SearchFields[k].SelNullText = "请选择"}
//		if mc.SearchFields[k].Width <=0 { mc.SearchFields[k].Width = 120}
//		if mc.SearchFields[k].Height <=0 { mc.SearchFields[k].Height = 60}
//		if mc.SearchFields[k].Values == nil { mc.SearchFields[k].Values = []string{"?"}}
//	}
//}
