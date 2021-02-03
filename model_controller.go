package mc

import "fmt"

type ModelController struct {
    Controller
}

//初始化控制器
func (ctrl *ModelController) Initialize(c *Context, childCtrl IController) {
    ctrl.Controller.Initialize(c, childCtrl)
    ctrl.ChildController.ModelInitializeBefore() //模型初始化之前
    //根据请求参数自动初始化model
    if c.ModelName == "" {
        c.ModelName = fmt.Sprintf("%s_%s", c.ModuleName, c.ControllerName)
    }
    ctrl.Model = NewModel(ctrl.Context.ModelName)
    ctrl.Assign.Model = ctrl.Model
    ctrl.ChildController.ModelInitializeAfter() //模型初始化之后
}

//首页操作
func (ctrl *ModelController) IndexAct() {
    if ctrl.Context.Request.Method == "GET" {
        ctrl.ChildController.ModelListUIRenderBefore()
        ctrl.Assign.Result["form_items"] = ctrl.Model.CreateSearchItems(nil)
        ctrl.AbortWithHtml(nil)
        return
    } else if ctrl.Context.Request.Method == "POST" {
        page, pageSize := ctrl.Context.getPage()
        //查询选项
        qo := &QueryOption{
            Values:   ctrl.UrlValueToRequestValue(ctrl.Context.Request.PostForm),
            Page:     page,
            PageSize: pageSize,
            Order:    ctrl.Context.getOrder(),
        }
        ctrl.ChildController.ModelFindBefore(qo) //数据查询之前回调
        //数据查询
        data, foot, total, err := ctrl.Model.Find(qo)
        if err != nil {
            panic(err)
        }
        //输出
        result := &Result{
            Data: data,
            ExtraData: map[string]interface{}{
                option.Response.TotalName: total,
                option.Response.FootName:  foot,
            },
        }
        ctrl.ChildController.ModelFindAfter(result) //数据查询之后回调
        ctrl.AbortWithMessage(result)
        return
    }
}

//添加操作
func (ctrl *ModelController) AddAct() {
    ctrl.Save(nil)
}

//编辑操作
func (ctrl *ModelController) EditAct() {
    pkValue := ctrl.Context.DefaultQuery(ctrl.Model.attr.Pk, "")
    if pkValue == "" {
        panic(&Result{Message: "请选择一条记录进行操作"})
    }
    ctrl.Save(pkValue)
}

//删除操作
func (ctrl *ModelController) DelAct() {
    ids := ctrl.Context.PostFormArray("id")
    if len(ids) <= 0 {
        panic(&Result{Message: "请选择至少一条记录进行操作"})
    }

    if _, err := ctrl.Model.Delete(ids); err != nil {
        panic(err)
    } else {
        result := &Result{Message: "数据删除成功"}
        ctrl.ChildController.ModelSaveAfter(result)
        ctrl.AbortWithMessage(result)
    }
}

// 仅供Add或Edit调用
func (ctrl *ModelController) Save(pkValue interface{}) {
    if ctrl.Context.Request.Method == "GET" { // GET 界面
        var rowData map[string]interface{}
        if pkValue != nil {
            qo := &QueryOption{
                ExtraWhere: []interface{}{fmt.Sprintf("%s = ?", ctrl.Model.FieldAddAlias(ctrl.Model.attr.Pk)), pkValue},
            }
            ctrl.ChildController.ModelEditTakeBefore(qo) ////编辑查询数据前回调
            if row, exist, err := ctrl.Model.TakeForEdit(qo); err != nil {
                panic(err)
            } else if !exist {
                panic(&Result{Message: "记录未找到"})
            } else {
                rowData = row
            }
        }
        ctrl.ChildController.ModelEditTakeAfter(rowData) //编辑查询数据后回调
        ctrl.Assign.Result["form_items"] = ctrl.Model.CreateEditItems(rowData)
        ctrl.AbortWithHtml(nil)
        return

    } else if ctrl.Context.Request.Method == "POST" { //POST 提交
        data := ctrl.UrlValueToRequestValue(ctrl.Context.Request.PostForm)
        ctrl.ChildController.ModelSaveBefore(data) //保存数据之前回调
        if _, err := ctrl.Model.Save(data, pkValue); err != nil {
            panic(err)
        } else {
            result := &Result{Message: "数据保存成功"}
            ctrl.ChildController.ModelSaveAfter(result)
            ctrl.AbortWithMessage(result)
        }
    }
}

