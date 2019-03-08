package admin

import (
	"github.com/gin-gonic/gin"
	"logic"
	"strconv"
	"strings"
)

// 网站信息配置
func (this *AdminControl) ConfigInfo (c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		_type, _ := strconv.Atoi(c.DefaultPostForm("type", "0"))
		if _type == 1 {

		} else if _type == 2 {

		} else {
			returnJson(c,1001,nil)
		}
		platform := c.DefaultPostForm("platform", "")
		record_no := c.DefaultPostForm("record_no", "")
		wx := c.DefaultPostForm("wx", "")
		tel := c.DefaultPostForm("tel", "")
		address := c.DefaultPostForm("address", "")
		

		iRelust := logic.GetSysMenu().GetMenuListByPage(int32(page), keyword)
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	returnHtml(c, "config-info.html", gin.H{
		"keyword": keyword,
	})
	return
}