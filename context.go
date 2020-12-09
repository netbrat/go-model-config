package mc

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

//路由请求参数
type Reqs struct {
	Module         string //模块
	RealModule     string //实际使用的模块
	Controller     string //控制器
	RealController string //实际使用的控制器
	Action         string //操作方法
	Model          string //模型
}

//分页信息
type PageInfo struct {
	Total    int
	Page     int
	PageSize int
	Offset   int
}

//上下文
type Context struct {
	Reqs     *Reqs
	PageInfo *PageInfo
	*gin.Context
}

func (c *Context) GetPage() (page int, pageSize int) {
	if c.Request.Method == "GET" {
		page = cast.ToInt(c.Query("page"))
		pageSize = cast.ToInt(c.Query("page_size"))
	} else {
		page = cast.ToInt(c.Request.Form["page"])
		page = cast.ToInt(c.Request.Form["page_size"])
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = option.PageSize
	}
	return
}

//上下文初始化
func (c *Context) Init() {
	//分析请求路径
	params := strings.Split(strings.ToLower(c.Request.URL.Path), "/")
	params = append(params, "", "", "", "")
	c.Reqs = &Reqs{
		Module:         params[1],
		RealModule:     params[1],
		Controller:     params[2],
		RealController: params[2],
		Action:         params[3],
		Model:          params[4],
	}

	//页码处理
	c.PageInfo = &PageInfo{}
	c.PageInfo.Page, c.PageInfo.PageSize = c.GetPage()
	c.PageInfo.Offset = (c.PageInfo.Page - 1) * c.PageInfo.PageSize

}
