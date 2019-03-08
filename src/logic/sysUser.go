package logic

import (
	"encoding/json"
	"utils"
)

type SysUser struct {
	Id       int32  `json:"id" bson:"id"`
	Nickname string `json:"nickname" bson:"nickname"`   // 昵称
	Password string `json:"password" bson:"password"`   // 密码
	Phone    string `json:"phone" bson:"phone"`         // 手机号
	Sex      int    `json:"sex" bson:"sex"`             // 性别 0: 女 1: 男
	RoleId   int    `json:"role_id" bson:"role_id"`     // 角色ID
	Remark   string `json:"remark" bson:"remark"`       // 备注
	RoleName string `json:"role_name" bson:"role_name"` // 角色名称
	Times    int32  `json:"times" bson:"times"`         // 登录次数
	Deleted  int    `json:"deleted" bson:"deleted"`     // 状态 0:正常 1: 冻结
	CreateAt int64  `json:"create_at" bson:"create_at"`
}

type sysUserList struct {
	List  []*SysUser `json:"list" bson:"list"`
	Total int64      `json:"total" bson:"total"`
}

// 分页获取管理员列表
func (this *SysUser) GetUserListByPage (page int32, keyword  string) interface{} {
	jsonStr, _ := db.Call("Proc_SysUser_pagination_v1.0", page, pageSize, keyword)
	info := &sysUserList{}
	json.Unmarshal([]byte(jsonStr), &info)
	iResult := make(map[string]interface{}, 0)
	iResult["total"] = info.Total
	iResult["pageSize"] = pageSize
	iResult["data"] = info.List
	return iResult
}

// 冻结管理员
func (this *SysUser) DelSysUserById (userId int64) int {
	jsonStr, _ := db.Call("Proc_SysUser_delById_v1.0", userId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询管理员详情
func (this *SysUser) GetUserInfoById (adminId int64) *SysUser {
	jsonStr, _ := db.Call("Proc_SysUser_infoById_v1.0", adminId)
	info := []*SysUser{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]
}

// 添加、编辑管理员
func (this *SysUser) ModifySysUser (nickname, phone, password, remark string, roleId, sex, adminId int64) int {
	if len(password) > 0  {
		password = utils.Md5(password)
	}
	jsonStr, _ := db.Call("Proc_SysUser_modify_v1.0", nickname, phone, sex, password, remark, roleId, adminId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}