package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"utils"
)

type ISysUser interface {
	GetSysUserListByPage(int32, string) interface{}
	DelSysUserById(int64) int
	GetUserInfoById(int64) *SysUser
	ModifySysUser(string, string, string, string, int64, int64, int64) int
	Login(string, string, *gin.Context) int
}

type SysUser struct {
	Id       int32
	Nickname string                                     // 昵称
	Password string                                     // 密码
	Phone    string                                     // 手机号
	Sex      int                                        // 性别 0: 女 1: 男
	RoleId   int `json:"role_id" bson:"role_id"`        // 角色ID
	Remark   string                                     // 备注
	RoleName string `json:"role_name" bson:"role_name"` // 角色名称
	Times    int32                                      // 登录次数
	Deleted  int                                        // 状态 0:正常 1: 冻结
	CreateAt int64 `json:"create_at" bson:"create_at"`
}

type sysUserList struct {
	List  []*SysUser
	Total int64
}

// 分页获取管理员列表
func (su *SysUser) GetSysUserListByPage (page int32, keyword string) interface{} {
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
func (su *SysUser) DelSysUserById (userId int64) int {
	jsonStr, _ := db.Call("Proc_SysUser_delById_v1.0", userId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询管理员详情
func (su *SysUser) GetUserInfoById (adminId int64) *SysUser {
	jsonStr, _ := db.Call("Proc_SysUser_infoById_v1.0", adminId)
	info := []*SysUser{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]
}

// 添加、编辑管理员
func (su *SysUser) ModifySysUser (nickname, phone, password, remark string, roleId, sex, adminId int64) int {
	if len(password) > 0  {
		password = utils.Md5(password)
	}
	jsonStr, _ := db.Call("Proc_SysUser_modify_v1.0", nickname, phone, sex, password, remark, roleId, adminId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

//后台登录
func (su *SysUser) Login (acconut, password string, c *gin.Context) int {
	jsonStr, _ := db.Call("Proc_SysUser_login_v1.0", acconut, password)
	info := []*SysUser{}
	json.Unmarshal([]byte(jsonStr), &info)
	if len(info) == 0 {
		return 1021
	}
	if info[0].Deleted == 1 {
		return 1020
	}
	//保存session
	sessionId, _ := c.Cookie("session_id")
	utils.GetSeesionMgr(c).Set(sessionId, "sysInfo", info[0])
	return 0
}


