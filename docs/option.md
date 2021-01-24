
# 配置项

### 设置配置项


```
# 获取gin的Engine
r := gin.Default() 

# 获取mc的默认配置
option := mc.Default(r)

# 设置mc的配置
option.ErrorTemplate = "base/error.html"
...

```

### mc配置项说明


配置项 | 子配置项 | 类型 | 默认值 | 说明
---|---|---|---|---
DefaultConnName | | string | default | 默认数据库连接名
ErrorTemplate | | string | error.html | 错误页面模版 
ModelConfigsFilePath | | string | ./mcconfigs/ | 自定义模型配置文件存放路径
Router | | | | 路由选项
. | UrlPathSep | string | / | URL路径之间的分割符号（不能使用_下线线）
. | UrlHtmlSuffix | string | html | URL伪静态后缀设置
. | ControllerMap | map[string]map[string]IController | | 控制器map
. | BaseModuleName | string | base | 全局基础模块key
. | BaseControllerName | string | base | 全局基础控制器key
. | ModuleBaseControllerName | string | base | 当前模块下基础控制器key
Response | | | | 响应项
. | CodeName | string | code | 代码项的key
. | MessageName | string | msg | 消息项的key
. | DataName | string | data | 数据项的key
. | TotalName | string | total | 总记录数或影响的记录数项的key
. | FootName | string | count | 表尾汇总数据项的key
. | SuccessCodeValue | string | 0000 | 成功代码值
. | FailCodeValue | string | 1000 | 失败默认代码值
. | AjaxRenderType | mc.RenderType | mc.RenderTypeJSON | 默认ajax渲染类型
Request | | | | 请求项
. | OrderName | string | order | 代码项的key
. | PageName | string | page | 消息项的key
. | PageSizeName | string | limit | 数据项的key
. | PageSizeValue | string | 50 | 总记录数或影响的记录数项的key
. | ContextInitializeStartFunc | func(c *gin.Context) (err error) |  | 初始化上下文前回调
. | ContextInitializeEndFunc | func(ctx *mc.Context) (err error) |  | 初始化上下文后回调
Auth | | | | 权限项
. | RowAuthModels | []string |  | 行权限model列表
. | GetAuthFunc | func() *mc.Auth |  | 获取权限回调