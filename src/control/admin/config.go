package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// 网站信息配置
func (this *AdminControl) ConfigInfo (c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		_type, _ := strconv.Atoi(c.DefaultPostForm("type", "0"))
		iRelust := 1003
		if _type == 1 {

		} else if _type == 2 {
			content := make(map[string]interface{})
			content["platform"] = c.DefaultPostForm("platform", "")
			content["record_no"] = c.DefaultPostForm("record_no", "")
			content["wx"] = c.DefaultPostForm("wx", "")
			content["tel"] = c.DefaultPostForm("tel", "")
			content["address"] = c.DefaultPostForm("address", "")
		} else {
			returnJson(c,1001,nil)
		}
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	returnHtml(c, "config-info.html", gin.H{
		"keyword": keyword,
	})
	return
}