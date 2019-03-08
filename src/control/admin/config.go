package admin

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-gonic/gin.v1/json"
	"logic"
	"strconv"
	"strings"
)

// 网站信息配置
func (this *AdminControl) ConfigInfo (c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		_type, _ := strconv.Atoi(c.DefaultPostForm("type", "1"))
		info := logic.GetConfig().GetConfigInfoByType(_type)
		returnJson(c, 0, gin.H{
			"content": info["content"],
			"keys": info["keys"],
		})
		return
	}
	returnHtml(c, "config-info.html", nil)
	return
}

// 添加、编辑网站信息
func (this *AdminControl) ConfigModify (c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		_type, _ := strconv.Atoi(c.DefaultPostForm("type", "0"))
		content := make(map[string]interface{})
		if _type == 1 {
			// 全局SEO 关键字、描述
			content = gin.H{
				"keyword":     c.DefaultPostForm("keyword", ""),
				"description": c.DefaultPostForm("description", ""),
			}
		} else if _type == 2 {
			// 网站配置
			content = gin.H{
				"platform":  c.DefaultPostForm("platform", ""),
				"record_no": c.DefaultPostForm("record_no", ""),
				"wx":        c.DefaultPostForm("wx", ""),
				"tel":       c.DefaultPostForm("tel", ""),
				"address":   c.DefaultPostForm("address", ""),
			}

		} else {
			returnJson(c,1001,nil)
		}
		jsonStr, _ := json.Marshal(content)
		iRelust := logic.GetConfig().ModifyConfig(_type, string(jsonStr))
		returnJson(c, 0, iRelust)
		return
	}
	returnHtml(c, "config-info.html", nil)
	return
}