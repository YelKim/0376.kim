package admin

import (
	"conf"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"utils"
)

type AdminControl struct{}

//上传状态映射表，国际化用户需考虑此处数据的国际化
var ueditorMap = gin.H{
	"ERROR_SIZE_EXCEED":      "文件大小超出网站限制",
	"ERROR_TYPE_NOT_ALLOWED": "文件类型不允许",
	"ERROR_FILE_MOVE":        "文件保存时出错",
	"ERROR_FILE_NOT_FOUND":   "找不到上传文件",
}

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
	"1014": "上传文件大小超出",
	"1015": "上传文件文件格式有误",
	"1016": "文件上传失败",
	"1017": "不能选自己作为父类",
	"1018": "验证码不能为空",
	"1019": "账号密码不能为空",
	"1020": "账号已被冻结",
	"1021": "账号密码错误",
	"1022": "验证码错误",
}

// 返回Json 字符串
func returnJson(c *gin.Context, code int, data interface{}) {
	iResult := gin.H{
		"code": code,
	}
	if data != nil {
		iResult["data"] = data
	}
	if code >= 0 {
		iResult["message"] = errMessage[strconv.Itoa(code)]
	}
	c.JSON(200, iResult)
}

// 返回Html模板 字符串
func returnHtml(c *gin.Context, views string, data interface{}) {
	c.HTML(200, views, data)
}

//上传文件
func (ac *AdminControl) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("upfile")
	if err != nil {
		returnJson(c, 1000, nil)
		return
	}
	//判断文件的大小
	if int32(header.Size) > conf.GetConfig().FileSize*1024*1024 {
		returnJson(c, 1014, nil)
		return
	}
	//判断文件格式
	i := true
	extArr := strings.Split(header.Header["Content-Type"][0], "/")
	for _, v := range conf.GetConfig().FileExt {
		if v == extArr[1] {
			i = false
			break
		}
	}
	if i {
		returnJson(c, 1015, nil)
		return
	}
	//判断文件夹是否存在
	filePath := conf.GetConfig().StaticDir + "/tmp/" + time.Now().Format("20060102")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, 0755)
	}

	md5 := []rune(utils.Md5(strconv.Itoa(int(time.Now().UnixNano()))))
	fileName := filePath + "/" + string(md5[8:24]) + "." + extArr[1]
	out, err := os.Create(fileName)
	if err != nil {
		returnJson(c, 1016, nil)
		return
	}
	defer out.Close()
	io.Copy(out, file)
	returnJson(c, 0, gin.H{"path": strings.Trim(fileName, ".")})
	return
}

// 百度编辑器配置、文件上传
func (ac *AdminControl) Ueditor(c *gin.Context) {
	action := c.DefaultQuery("action", "config")
	switch action {
	case "config":
		file, err := ioutil.ReadFile("./lib/ueditor/1.4.3/config.json")
		if err != nil {
			return
		}
		c.Header("Content-type", "application/json")
		c.Writer.Write(file)
		return
		break
	case "image":
		file, header, err := c.Request.FormFile("upfile")
		if err != nil {
			c.JSON(200, gin.H{"tate": ueditorMap["ERROR_FILE_NOT_FOUND"]})
			return
		}
		//判断文件的大小
		if int32(header.Size) > conf.GetConfig().FileSize*1024*1024 {
			c.JSON(200, gin.H{"tate": ueditorMap["ERROR_SIZE_EXCEED"]})
			return
		}
		//判断文件格式
		i := true
		extArr := strings.Split(header.Header["Content-Type"][0], "/")
		for _, v := range conf.GetConfig().FileExt {
			if v == extArr[1] {
				i = false
				break
			}
		}
		if i {
			c.JSON(200, gin.H{"tate": ueditorMap["ERROR_TYPE_NOT_ALLOWED"]})
			return
		}
		//判断文件夹是否存在
		filePath := conf.GetConfig().StaticDir + "/tmp/" + time.Now().Format("20060102")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			os.MkdirAll(filePath, 0755)
		}
		md5 := []rune(utils.Md5(strconv.Itoa(int(time.Now().UnixNano()))))
		fileName := filePath + "/" + string(md5[8:24]) + "." + extArr[1]
		out, err := os.Create(fileName)
		if err != nil {
			c.JSON(200, gin.H{"tate": ueditorMap["ERROR_FILE_MOVE"]})
			return
		}
		defer out.Close()
		io.Copy(out, file)
		c.JSON(200, gin.H{
			"state":    "SUCCESS",
			"url":      strings.Trim(fileName, "."),
			"title":    header.Filename,
			"original": header.Filename,
			"type":     extArr[1],
			"size":     header.Size,
		})
		return
	}
}
