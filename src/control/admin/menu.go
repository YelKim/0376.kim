package admin

import (
	"github.com/gin-gonic/gin"
	"logic"
	"strconv"
	"strings"
)

// 栏目列表
func (this *AdminControl) MenuList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		keyword := c.DefaultPostForm("keyword", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))
		iRelust := logic.GetMenu().GetMenuListByPage(int32(page), keyword)
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	returnHtml(c, "menu.html", gin.H{
		"keyword": keyword,
	})
	return
}

// 子栏目列表
func (this *AdminControl) MenuChildList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
		if parentId == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.GetMenu().GetMenuChildListByParentId(int64(parentId))
		returnJson(c, 0, iRelust)
		return
	}
	returnJson(c, 1000, nil)
	return
}

// 删除栏目
func (this *AdminControl) MenuDel(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		menuId, _ := strconv.Atoi(c.DefaultPostForm("menu_id", "0"))
		if menuId == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.GetMenu().DelMenuById(int64(menuId))
		returnJson(c, iRelust, nil)
		return
	}
	returnJson(c, 1000, nil)
	return
}

// 编辑栏目
func (this *AdminControl) MenuModify(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		menuId, _ := strconv.Atoi(c.DefaultPostForm("menu_id", "0"))
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
		sort, _ := strconv.Atoi(c.DefaultPostForm("sort", "999"))
		title := c.DefaultPostForm("title", "")
		name := c.DefaultPostForm("name", "")
		keyword := c.DefaultPostForm("keyword", "")
		description := c.DefaultPostForm("description", "")
		url := c.DefaultPostForm("url", "")
		// 2、3级菜单 control action 不能为空
		if len(name) == 0 || len(url) == 0 {
			returnJson(c, 1001, nil)
		}
		iResult := logic.GetMenu().ModifyMenu(name, title, keyword, description, url, int64(parentId), int64(sort), int64(menuId))
		returnJson(c, iResult, nil)
		return
	}
	// 获取栏目详情
	menuId, _ := strconv.Atoi(c.DefaultQuery("menu_id", "0"))
	info := &logic.Menu{}
	info.Sort = 999
	// 获取栏目分类列表
	if menuId > 0 {
		info = logic.GetMenu().GetMenuInfoById(int64(menuId))
	}
	menuTree := logic.GetMenu().GetMenuListByLevel(1)
	returnHtml(c, "menu-modify.html", gin.H{
		"info":     info,
		"menuTree": menuTree,
	})
	return
}
