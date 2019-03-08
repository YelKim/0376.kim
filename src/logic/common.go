package logic

import (
	"conf"
	"sync"
	"utils/mysql"
)

var db = mysql.GetMysql()
var once sync.Once
var pageSize = conf.GetConfig().PageSize

// 后台菜单
var sysMenu *SysMenu
func GetSysMenu() *SysMenu {
	once.Do(func() {
		sysMenu = &SysMenu{}
	})
	return sysMenu
}

// 管理员角色
var sysRole *SysRole
func GetSysRole() *SysRole {
	once.Do(func() {
		sysRole = &SysRole{}
	})
	return sysRole
}

// 管理员
var sysUser *SysUser
func GetSysUser() *SysUser {
	once.Do(func() {
		sysUser = &SysUser{}
	})
	return sysUser
}

