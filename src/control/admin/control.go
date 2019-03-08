package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type AdminControl struct {}

// 错误码
var errMessage = gin.H{
	"0":    "操作成功",
	"1000": "非法操作",
	"1001": "参数错误",
	"1002": "存在子菜单, 请先处理子菜单",
	"1003": "操作失败",
	"1004": "名称或标题已存在",
	"1005": "手机号已存在",
	"1006": "账号或者密码错误",
	"1007": "验证码错误",
	"1008": "找不到上传的临时文件",
	"1009": "文件没有通过 'HTTP POST' 上传",
	"1010": "文件保存出错",
	"1011": "文件格式错误",
	"1012": "文件大小是否超出限制",
	"1013": "结束时间大于开始时间",
}

// 返回Json 字符串
func returnJson(c *gin.Context, code int, data interface{}) {
	iResult := gin.H{
		"code": code,
	}
	if data != nil {
		iResult["data"] = data
	}
	if code > 0 {
		iResult["message"] = errMessage[strconv.Itoa(code)]
	}
	c.JSON(200, iResult)
}

// 返回Html模板 字符串
func returnHtml(c *gin.Context, views string, data interface{}) {
	c.HTML(200, views, data)
}