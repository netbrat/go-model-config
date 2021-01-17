package mc

import "fmt"

type ModelController struct {
    Controller
    Model *Model
}


func (ctrl *ModelController) Initialize(c *Context) {
    if c.ModelName == "" {
        c.ModelName = fmt.Sprintf("%s_%s", c.ModuleName,c.ControllerName)
    }
    ctrl.Controller.Initialize(c)
    ctrl.Model = NewModel(ctrl.Context.ModelName)
    ctrl.Assign.Model = ctrl.Model
}


func (ctrl *ModelController) Index(){
    if ctrl.Context.Request.Method=="GET"{
        ctrl.Model.CreateSearchItems(nil)
        ctrl.AbortWithSuccess(Result{RenderType:RenderTypeHTML})
    }else if ctrl.Context.Request.Method == "POST"{
        page, pageSize := ctrl.Context.getPage()
        //查询选项
        qo := &QueryOption{
            Values: ctrl.UrlValueToRequestValue(ctrl.Context.Request.PostForm),
            Page: page,
            PageSize:pageSize,
            Order:ctrl.Context.getOrder(),
        }
        fmt.Println(qo.Values)
        //数据查询
        data, footer, total, err := ctrl.Model.Find(qo)
        if err != nil{
            ctrl.AbortWithError(Result{Message:err.Error()})
            return
        }
        //输出
        result := Result{
            Data: data,
            ExtraData: map[string]interface{}{
                option.Response.TotalName:  total,
                option.Response.FooterName: footer,
            },
        }
        ctrl.AbortWithSuccess(result)
    }
}


func (ctrl *ModelController) Add(){
    ctrl.Save(nil)
}

// 编辑操作
func (ctrl *ModelController) Edit(){
    pkValue := ctrl.Context.DefaultQuery(ctrl.Model.attr.Pk,"")
    if pkValue == ""{
        ctrl.AbortWithError(Result{Message:"请选择一条记录进行操作"})
        return
    }
    ctrl.Save(pkValue)
}

// 保存数据（不用于界面操作），仅供Add或Edit调用
func (ctrl *ModelController) Save(pkValue interface{}){
    if ctrl.Context.Request.Method == "GET" { // GET 界面
        var rowData map[string]interface{}
        if pkValue != nil {
            qo := &QueryOption{
                ExtraWhere:  []interface{}{fmt.Sprintf("%s = ?", ctrl.Model.attr.Pk), pkValue},
                NotSearch:   true,
                NotColAuth: true,
            }
            if row, exist, err := ctrl.Model.Take(qo); err != nil{
                ctrl.AbortWithError(Result{Message:err.Error()})
                return
            }else if !exist {
                ctrl.AbortWithError(Result{Message:"记录未找到"})
                return
            } else{
                rowData = row
            }
        }
        ctrl.Model.CreateEditItems(rowData)
        ctrl.AbortWithSuccess(Result{RenderType:RenderTypeHTML})

    } else if ctrl.Context.Request.Method == "POST" { //POST 提交

        data := ctrl.UrlValueToRequestValue(ctrl.Context.Request.PostForm)
        if _, err := ctrl.Model.Save(data, pkValue); err != nil{
            ctrl.AbortWithError(Result{Message:err.Error()})
        }else{
            ctrl.AbortWithSuccess(Result{Message:"数据保存成功"})
        }
    }
}

func (ctrl *ModelController) Del() {
    ids := ctrl.Context.PostFormArray("id")
    if len(ids)<=0 {
        ctrl.AbortWithError(Result{Message:"请选择至少一条记录进行操作"})
        return
    }
    if _, err := ctrl.Model.Delete(ids); err != nil{
        ctrl.AbortWithError(Result{Message:err.Error()})
    }else{
        ctrl.AbortWithSuccess(Result{Message:"数据删除成功"})
    }
}