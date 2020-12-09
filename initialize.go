package mc

type Option struct {
	PageSize        int       //默认页记录数
	ErrorTemplate   string    //错误页面模版
	RouterMap       RouterMap //路由列表
	NotAuthRedirect string    //未登录跳转到页面地址
}

var option = &Option{}

func Initialize(opt *Option) {
	option = opt
}
