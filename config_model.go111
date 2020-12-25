package mc

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"strings"
)

type ConfigModel struct{
	*BaseModel
	Config *Config
	Auth	interface{}
}


func (cm *ConfigModel) SetConfig(configName string)(err error){
	if cm.Config, err = GetConfig(configName); err != nil{
		return err
	}
	cm.SetDbOption(cm.Config.ConnName, cm.Config.DbName, cm.Config.Table, cm.Config.Pk, cm.Config.AutoIncrement, cm.Config.UniqueFields)
	return
}


func (cm *ConfigModel) GetKv(kvName string, so SearchOption) (data []map[string]interface{}, err error){
	if kvName == "" { kvName = "default"}
	if !InArray(kvName,cm.Config.Kvs){
		return nil, fmt.Errorf("配置中不存在 [%s] kv 项配置", kvName)
	}
	db := cm.parseWhere(GetDB(cm.Config.ConnName), so)
	fields := cm.parseKvFields(kvName, so)

	return
}


// @title parseKvFields
// @description 分析KV查询条件
// @param 	so 		SearchOption	模型查询字段各项的值
// @return	fields	[]string		最终需要查询的KV字段名数组
func (cm *ConfigModel) parseKvFields(kvName string, so SearchOption) (fields []string){
	fields = make([]string, 0)

	// kv配置中的字段
	kv, ok := Kv{}, false
	if kv, ok = cm.Config.Kvs[kvName]; !ok{
		return
	}
	keySep := fmt.Sprintf("'%s',%s.", kv.KeySep, cm.Config.BaseSearch.Alias)
	valueSep := fmt.Sprintf("'%s',%s.", kv.ValueSep, cm.Config.BaseSearch.Alias)
	keyField := fmt.Sprintf("CONCAT(%s.%s) AS _key", cm.Config.BaseSearch.Alias, strings.Join(kv.KeyFields, keySep))
	ValueField := fmt.Sprintf("CONCAT(%s.%s) AS _value", cm.Config.BaseSearch.Alias, strings.Join(kv.ValueFields, valueSep))
	fields = append(fields, keyField, ValueField)

	// 树型必备字段
	if cm.Config.IsTree {
		fields = append(fields, cm.Config.TreePathField, cm.Config.TreeLevelField)
	}
	// 附加字段
	if so.Fields != nil {
		fields = append(fields, so.Fields...)
	}
	return
}


// @title parseWhere
// @description 分析查询条件
// @param	db		*gorm.DB
// @param 	so 		SearchOption	模型查询字段各项的值
// @return	rdb		*gorm.DB		附带最终的查询条件的db对象
func (cm *ConfigModel) parseWhere(db *gorm.DB, so SearchOption) (rdb *gorm.DB) {
	rdb = db
	// 模型全局查询条件
	if cm.Config.BaseSearch.Where != ""{
		rdb = rdb.Where(cm.Config.BaseSearch.Where)
	}
	// so中额外附加的查询条件
	rdb = cm.BaseModel.parseWhere(rdb, so)

	// 模型各查询字段
	if so.IsSearch {
		for _, f := range cm.Config.SearchFields {
			// 该查询字段未带条件配置，跳过
			if f.Where == "" {
				continue
			}
			// 未传入查询值时，使用默认值
			if cast.ToString(so.SearchValues[f.Name]) == "" {
				if f.Default != nil {
					delete(so.SearchValues, f.Name)
				} else {
					so.SearchValues[f.Name] = f.Default
				}
			}
			// 查询值与查询条件匹配
			values := make([]interface{}, 0)
			for _, v := range f.Values {
				if v == "?" {
					values = append(values, so.SearchValues[f.Name])
				} else {
					values = append(values, strings.ReplaceAll(v, "?", cast.ToString(so.SearchValues[f.Name])))
				}
			}
			rdb = rdb.Where(f.Where, values...)
		}
	}

	// 不使用级权限
	if so.NotRowAuth { return }

	return

}
