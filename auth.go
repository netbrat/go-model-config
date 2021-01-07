package mc

//权限
type Auth struct {
    User	map[string]interface{}		//用户信息
    UserAuth	map[string]interface{}	//用户权限
    RoleAuth	map[string]interface{}  //角色权限
}


