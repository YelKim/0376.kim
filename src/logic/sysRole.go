package logic

import (
	"encoding/json"
)

type ISysRole interface {
	GetRoleListByPage(int32, string) interface{}
	GetRoleInfoById(int64) *SysRole
	ModifySysRole(int64, string, string, string) int
	GetSysRoleMenuListByRoleId(rolId int64) []int32
	GetSysRoleList() []*SysRole
}

type SysRole struct {
	Id      int32  `json:"id" bson:"id"`
	Title   string `json:"title" bson:"title"`   //角色名称
	Remark  string `json:"remark" bson:"remark"` //备注
	UpateAt int64  `json:"update_at" bson:"update_at"`
}

type sysRoleList struct {
	List  []*SysRole `json:"list" bson:"list"`
	Total int64      `json:"total" bson:"total"`
}

// 分页获取管理员角色列表
func (sr *SysRole) GetRoleListByPage (page int32, keyword string) interface{} {
	jsonStr, _ := db.Call("Proc_SysRole_pagination_v1.0", page, pageSize, keyword)
	info := &sysRoleList{}
	json.Unmarshal([]byte(jsonStr), &info)
	iResult := make(map[string]interface{}, 0)
	iResult["total"] = info.Total
	iResult["pageSize"] = pageSize
	iResult["data"] = info.List
	return iResult
}

// 根据ID查询管理员角色详情
func (sr *SysRole) GetRoleInfoById (roleId int64) *SysRole {
	jsonStr, _ := db.Call("Proc_SysRole_infoById_v1.0", roleId)
	info := []*SysRole{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]
}

// 添加、编辑管理员角色
func (sr *SysRole) ModifySysRole (roleId int64, title, remark, menuArr string) int {
	jsonStr, _ := db.Call("Proc_SysRole_modify_v1.0", title, remark, menuArr, roleId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据管理员角色ID获取角色对应的菜单id列表
func (sr *SysRole) GetSysRoleMenuListByRoleId (roleId int64) []int32 {
	jsonStr, _ := db.Call("Proc_SysRoleMenu_listByRoleId_v1.0", roleId)
	info := []map[string]int32{}
	json.Unmarshal([]byte(jsonStr), &info)
	iResult := []int32{}
	if len(info) > 0 {
		for _, v := range info {
			iResult = append(iResult, v["menu_id"])
		}
	}
	return iResult
}

// 获取所有角色，用于下列框显示
func (sr *SysRole) GetSysRoleList() []*SysRole {
	jsonStr, _ := db.Call("Proc_SysRole_list_v1.0")
	info := []*SysRole{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info
}