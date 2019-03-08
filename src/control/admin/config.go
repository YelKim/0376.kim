package admin

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-gonic/gin.v1/json"
	"logic"
	"strconv"
	"strings"
)

// 网站信息配置
func (this *AdminControl) ConfigInfo(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		_type, _ := strconv.Atoi(c.DefaultPostForm("type", "1"))
		info := logic.GetConfig().GetConfigInfoByType(_type)
		returnJson(c, 0, gin.H{
			"content": info["content"],
			"keys":    info["keys"],
		})
		return
	}
	returnHtml(c, "config-info.html", nil)
	return
}

// 编辑网站配置信息
func (this *AdminControl) ConfigModify(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		_type, _ := strconv.Atoi(c.DefaultPostForm("type", "0"))
		var content interface{}
		if _type == 2 {
			// 全局SEO 关键字、描述
			content = &logic.SeoConfig{
				Keyword:     c.DefaultPostForm("keyword", ""),
				Description: c.DefaultPostForm("description", ""),
			}
		} else if _type == 1 {
			// 网站配置
			content = logic.InfoConfig{
				Platform: c.DefaultPostForm("platform", ""),
				RecordNo: c.DefaultPostForm("record_no", ""),
				Wx:       c.DefaultPostForm("wx", ""),
				Tel:      c.DefaultPostForm("tel", ""),
				Address:  c.DefaultPostForm("address", ""),
			}
		} else {
			returnJson(c, 1001, nil)
		}
		jsonStr, _ := json.Marshal(content)
		iRelust := logic.GetConfig().ModifyConfig(_type, string(jsonStr))
		returnJson(c, iRelust, nil)
		return
	}
	returnHtml(c, "config-info.html", nil)
	return
}
