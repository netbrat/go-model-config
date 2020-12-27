package mc
//
//import (
//	"errors"
//	"fmt"
//	"github.com/spf13/cast"
//	"gorm.io/gorm"
//	"strings"
//)
//
//type KvsSearchOption struct {
//	//DB			*gorm.DB		//使用的数据查询对象，如为nil时，使用model内置db对象
//	KvName		string			//kv配置项名
//	ExtraWhere	interface{}		//附加的查询条件
//	returnPath	bool			//当模型为树型结构时，返回的key是否使用path代替
//	NotRowAuth	bool			//是否不使用行级权限过滤条件
//	Indent		string			//当模型为树型结构时，层级缩进符, 空字符串时不缩进
//	ExtraFields	[]string		//额外附加的查询字段
//}
//
//
//type configModel struct {
//	*model
//	config 			*Config
//	configName		string
//}
//
//
//// 新建一个自定义配制模型
//// @param configName string 配制名
//func NewConfigModel(configName string) (*configModel, error) {
//	//读取模型配置文件
//	config, err := GetConfig(configName)
//	if err != nil {
//		return nil, err
//	}
//	m, err := NewModel(config.ConnName, config.DbName, config.Table, config.BaseSearch.Alias, config.PkField, config.AutoIncrement, config.UniqueFields)
//	if err != nil {
//		return nil, err
//	}
//	cm := &configModel{
//		model:     m,
//		config:     config,
//		configName: configName,
//	}
//	return cm, nil
//}
//
//
//// 获取自自定义配制信息
//func (m *configModel) Config() Config {
//	return *m.config
//}
//
//
//// 获取自定义配制名
//func (m *configModel) ConfigName() string {
//	return m.configName
//}
//
//
//// 获取Kv键值列表
//func (m *configModel) GetKvs(so *KvsSearchOption) (result []map[string]interface{}, err error){
//	if so, err = m.checkKvsSearchOption(so); err != nil{
//		return
//	}
//	//db := NewDBSession(m.db)
//	m.parseWhere(m.db, so.ExtraWhere, nil, true, false)
//	fields := m.parseKvFields(so.KvName, so.ExtraFields)
//	var data []map[string]interface{}
//	m.db.Where("1=1")
//	m.db.Where("2=2")
//	if db := m.db.Select(fields).Find(&data); errors.Is(db.Error, gorm.ErrRecordNotFound) {
//		err = db.Error
//	}
//
//	return
//}
//
//
//
//// 分析查询条件（仅对通过NewConfigModel创建的模型有效）
//// @param db		这些条件作用于哪个数据连接对象上
//// @param extraWhere 额外的查询条件
//// @param searchValues 查询字段值
//// @param notSearch 是否使用查询字段条件
//// @param notRowAuth 是否使用行级权限进行过滤
//func (m *configModel) parseWhere(db *gorm.DB, extraWhere interface{}, searchValues map[string]interface{}, notSearch bool, notRowAuth bool){
//	//额外的查询条件
//	if extraWhere != nil {
//		db.Where(extraWhere)
//	}
//
//	//如果自定义模型未定义
//	if m.config == nil{
//		return
//	}
//
//	// 模型全局查询条件
//	if m.config.BaseSearch.Where != "" {
//		db.Where(m.config.BaseSearch.Where)
//	}
//
//	// 模型各查询字段
//	if !notSearch{
//		for _, f := range m.config.SearchFields {
//			// 该查询字段未带条件配置，跳过
//			if f.Where == "" {
//				continue
//			}
//			// 未传入查询值时，使用默认值
//			if cast.ToString(searchValues[f.Name]) == "" {
//				if f.Default != nil {
//					delete(searchValues, f.Name)
//				} else {
//					searchValues[f.Name] = f.Default
//				}
//			}
//			// 查询值与查询条件匹配
//			values := make([]interface{}, 0)
//			for _, v := range f.Values {
//				if v == "?" {
//					values = append(values, searchValues[f.Name])
//				} else {
//					values = append(values, strings.ReplaceAll(v, "?", cast.ToString(searchValues[f.Name])))
//				}
//			}
//			db.Where(f.Where, values...)
//		}
//	}
//	if !notRowAuth {
//
//	}
//}
//
//
//// 分析kv字段数组 （仅对通过NewConfigModel创建的模型有效）
//// @param 	kvName  kv配置项名
//// @return	fields	[]string		最终需要查询的KV字段名数组
//func (m *configModel) parseKvFields(kvName string, extraFields []string) (fields []string){
//	fields = make([]string, 0)
//
//	// kv配置中的字段
//	kv, ok := ConfKv{}, false
//	if kv, ok = m.config.Kvs[kvName]; !ok{
//		return
//	}
//	keySep := fmt.Sprintf("'%s',%s.", kv.KeySep, m.config.BaseSearch.Alias)
//	valueSep := fmt.Sprintf("'%s',%s.", kv.ValueSep, m.config.BaseSearch.Alias)
//	keyField := fmt.Sprintf("CONCAT(%s.%s) AS a_key", m.config.BaseSearch.Alias, strings.Join(kv.KeyFields, keySep))
//	ValueField := fmt.Sprintf("CONCAT(%s.%s) AS a_value", m.config.BaseSearch.Alias, strings.Join(kv.ValueFields, valueSep))
//	fields = append(fields, keyField, ValueField)
//
//	// 树型必备字段
//	if m.config.IsTree {
//		fields = append(fields, m.config.TreePathField, m.config.TreeLevelField)
//	}
//	// 附加字段
//	if extraFields != nil {
//		fields = append(fields, extraFields...)
//	}
//	return
//}
//
//
//// 检查kv查询选项
//func (m *configModel) checkKvsSearchOption(so *KvsSearchOption) (rso *KvsSearchOption, err error){
//	rso = so
//	if rso == nil {
//		rso = &KvsSearchOption{KvName: "default"}
//	}
//	if rso.KvName == "" {
//		rso.KvName = "default"
//	}
//	if !InArray(rso.KvName, m.config.Kvs){
//		err = fmt.Errorf("配置中不存在 [%s] kv 项配置", rso.KvName)
//	}
//	return
//}