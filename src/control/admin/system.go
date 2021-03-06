package admin

import (
	"github.com/gin-gonic/gin"
	"logic"
	"strconv"
	"strings"
)

var sysmenu *logic.SysMenu
var sysuser *logic.SysUser
var sysrole *logic.SysRole

// 后台菜单列表
func (ac *AdminControl) SysMenuList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		keyword := c.DefaultPostForm("keyword", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))
		iRelust := logic.ISysMenu(sysmenu).GetSysMenuListByPage(int32(page), keyword)
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	returnHtml(c, "sys-menu.html", gin.H{
		"keyword": keyword,
	})
	return
}

// 后台子菜单列表
func (ac *AdminControl) SysMenuChildList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
		if parentId == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.ISysMenu(sysmenu).GetSysMenuChildListByParentId(int64(parentId))
		returnJson(c, 0, iRelust)
		return
	}
	returnJson(c, 1000, nil)
	return
}

// 删除后台菜单
func (ac *AdminControl) SysMenuDel(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		menuId, _ := strconv.Atoi(c.DefaultPostForm("menu_id", "0"))
		if menuId == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.ISysMenu(sysmenu).DelSysMenuById(int64(menuId))
		returnJson(c, iRelust, nil)
		return
	}
	returnJson(c, 1000, nil)
	return
}

// 编辑后台菜单
func (ac *AdminControl) SysMenuModify(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		menuId, _ := strconv.Atoi(c.DefaultPostForm("menu_id", "0"))
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
		sort, _ := strconv.Atoi(c.DefaultPostForm("sort", "999"))
		title := c.DefaultPostForm("title", "")
		icon := c.DefaultPostForm("icon", "")
		_control := c.DefaultPostForm("control", "")
		action := c.DefaultPostForm("action", "")
		// 2、3级菜单 control action 不能为空
		if parentId > 0 && (len(_control) == 0 || len(action) == 0) {
			returnJson(c, 1001, nil)
		}
		iResult := logic.ISysMenu(sysmenu).ModifySysMenu(int64(menuId), int64(parentId), int64(sort), title, icon, _control, action)
		returnJson(c, iResult, nil)
		return
	}
	// 获取后台菜单详情
	menuId, _ := strconv.Atoi(c.DefaultQuery("menu_id", "0"))
	info := &logic.SysMenu{}
	info.Sort = 999
	// 获取菜单分类列表
	if menuId > 0 {
		info = logic.ISysMenu(sysmenu).GetSysMenuInfoById(int64(menuId))
	}
	menuTree := logic.ISysMenu(sysmenu).GetSysMenuListByLevel(2)
	returnHtml(c, "sys-menu-modify.html", gin.H{
		"info":     info,
		"menuTree": menuTree,
	})
	return
}

// 管理员角色列表
func (ac *AdminControl) SysRoleList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		keyword := c.DefaultPostForm("keyword", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))
		iRelust := logic.ISysRole(sysrole).GetRoleListByPage(int32(page), keyword)
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	returnHtml(c, "sys-role.html", gin.H{
		"keyword": keyword,
	})
	return
}

// 编辑管理员角色
func (ac *AdminControl) SysRoleModify(c *gin.Context) {
	// ajax 提交
	if strings.ToUpper(c.Request.Method) == "POST" {
		roleId, _ := strconv.Atoi(c.DefaultPostForm("role_id", "0"))
		title := c.DefaultPostForm("title", "")
		remark := c.DefaultPostForm("remark", "")
		menuArr, _ := c.GetPostFormArray("menu_arr")
		if len(menuArr) == 0 || len(title) == 0 {
			returnJson(c, 1001, nil)
		}
		iResult := logic.ISysRole(sysrole).ModifySysRole(int64(roleId), title, remark, strings.Join(menuArr, ","))
		returnJson(c, iResult, nil)
		return
	}
	// 获取管理员角色详情
	roleId, _ := strconv.Atoi(c.DefaultQuery("role_id", "0"))
	info := &logic.SysRole{}
	// 获取角色ID对应的菜单列表ID
	menuRoleList := []int32{}
	if roleId > 0 {
		info = logic.ISysRole(sysrole).GetRoleInfoById(int64(roleId))
		menuRoleList = logic.ISysRole(sysrole).GetSysRoleMenuListByRoleId(int64(roleId))
	}
	// 获取菜单分类列表
	menuTree := logic.ISysMenu(sysmenu).GetSysMenuListByLevel(3)
	returnHtml(c, "sys-role-modify.html", gin.H{
		"info":         info,
		"menuTree":     menuTree,
		"menuRoleList": menuRoleList,
	})
	return
}

// 管理员列表
func (ac *AdminControl) SysUserList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		keyword := c.DefaultPostForm("keyword", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))
		iRelust := logic.ISysUser(sysuser).GetSysUserListByPage(int32(page), keyword)
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	returnHtml(c, "sys-user.html", gin.H{
		"keyword": keyword,
	})
	return
}

// 冻结、解冻管理员
func (ac *AdminControl) SysUserDel(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		adminId, _ := strconv.Atoi(c.DefaultPostForm("admin_id", "0"))
		if adminId == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.ISysUser(sysuser).DelSysUserById(int64(adminId))
		returnJson(c, iRelust, nil)
		return
	}
	returnJson(c, 1000, nil)
	return
}

// 编辑管理员
func (ac *AdminControl) SysUserModify(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		adminId, _ := strconv.Atoi(c.DefaultPostForm("admin_id", "0"))
		nickname := c.DefaultPostForm("nickname", "")
		phone := c.DefaultPostForm("phone", "")
		password := c.DefaultPostForm("password", "")
		remark := c.DefaultPostForm("remark", "")
		sex, _ := strconv.Atoi(c.DefaultPostForm("sex", "0"))
		roleId, _ := strconv.Atoi(c.DefaultPostForm("role_id", "0"))
		if roleId == 0 || len(phone) == 0 || len(nickname) == 0 || (adminId == 0 && len(password) == 0) {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.ISysUser(sysuser).ModifySysUser(nickname, phone, password, remark, int64(roleId), int64(sex), int64(adminId))
		returnJson(c, iRelust, nil)
		return
	}
	// 获取管理员详情
	adminId, _ := strconv.Atoi(c.DefaultQuery("admin_id", "0"))
	info := &logic.SysUser{}
	if adminId > 0 {
		info = logic.ISysUser(sysuser).GetUserInfoById(int64(adminId))
	}
	// 获取角色列表
	roleList := logic.ISysRole(sysrole).GetSysRoleList()
	returnHtml(c, "sys-user-modify.html", gin.H{
		"info":     info,
		"roleList": roleList,
	})
	return
}
