package mc


//权限对象
type Auth struct {
    isSuper bool                    //是否超级权限
    rowsAuth map[string][]string    //行权限map {modelName1:[value1,value2,...],modelName2:[*]} ,如果全部权限，则为[*]
    colsAuth map[string][]string    //列权限map {modelName1:[fieldName1, fieldName2,...],modelName2:[*]} ,如果全部权限，则为[*]
}

//是否超级管理员
func (auth *Auth) IsSuper() bool{
    return auth.isSuper
}

//设置权限
func (auth *Auth) Set(isSuper bool, rowsAuth map[string][]string, colsAuth map[string][]string){
    auth.isSuper = isSuper
    auth.rowsAuth = rowsAuth
    auth.colsAuth = colsAuth
}


//获取权限代码列表
func (auth *Auth) getAuth(modelName string, authType string) (authCodes []string, isAllAuth bool){
    authCodes = make([]string, 0)
    if auth.isSuper{
        isAllAuth = true
        return
    }
    if auth.rowsAuth == nil || auth.rowsAuth[modelName] == nil{
        return
    }
    if authType == "row"{
        authCodes = auth.rowsAuth[modelName]
    }else{
        authCodes = auth.colsAuth[modelName]
    }
    if authCodes[0] == "*" {
        isAllAuth = true
    }
    return
}

//检查是否有该权限
func (auth *Auth) checkAuth(modelName string, code string, authType string) bool{
    authCodes, isAllAuth := auth.getAuth(modelName, authType)
    if isAllAuth{
        return true
    }else if InArray(code, authCodes){
        return true
    }
    return false
}

//获取行权限
func (auth *Auth) GetRowAuth(modelName string) (rowAuth []string, isAllAuth bool){
    return auth.getAuth(modelName, "row")
}

//判断某列是否有权限
func (auth *Auth) CheckRowAuth(modelName string, value string) bool{
    return auth.checkAuth(modelName, value, "row")
}

//获取列权限
func (auth *Auth) GetColAuth(modelName string) (colAuth []string, isAllAuth bool){
    return auth.getAuth(modelName, "col")
}

//判断某列是否有权限
func (auth *Auth) CheckColAuth(modelName string, fieldName string) bool{
    return auth.checkAuth(modelName, fieldName, "col")
}






