package admin

import (
	"conf"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"logic"
	"os"
	"runtime"
	"strings"
	"utils"
)

//首页
func (this *AdminControl) Index(c *gin.Context) {
	// 获取后台左侧菜单
	menuTree := logic.GetSysMenu().GetSysMenuListByLevel(2)
	sysInfo := utils.GetSeesionMgr(c).Get("sysInfo")
	returnHtml(c, "index.html", gin.H{
		"menuTree": menuTree,
		"info": sysInfo,
	})
	return
}

//欢迎页
func (this *AdminControl) Welcome(c *gin.Context) {
	//服务器名称
	host, _ := os.Hostname()
	path, _ := os.Getwd()
	returnHtml(c, "welcome.html", gin.H{
		"host":    host,
		"path":    path,
		"os":      runtime.GOOS,
		"cpu":     runtime.NumCPU(),
		"version": runtime.Version(),
		"domain":  c.GetHeader("Host"),
		"port":    conf.GetConfig().AdminPort,
	})
	return
}

//登录
func (this *AdminControl) Login(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		code := c.DefaultPostForm("captcha", "")
		account := c.DefaultPostForm("account", "")
		password := c.DefaultPostForm("password", "")
		reald := utils.GetSeesionMgr(c).Get("captcha")
		if len(code) == 0{
			returnJson(c, 1018, nil)
			return
		}
		if len(password) == 0 || len(account) == 0 {
			returnJson(c, 1019, nil)
			return
		}
		// 验证验证码
		if  reald == nil || !captcha.VerifyString(reald.(string), code) {
			returnJson(c, 1022, nil)
			return
		}
		iResult := logic.GetSysUser().Login(account, utils.Md5(strings.Trim(password, "")), c)
		returnJson(c, iResult, nil)
		return
	}
	returnHtml(c, "login.html", nil)
	return
}

//验证码
func (this *AdminControl) Captcha(c *gin.Context) {
	code := captcha.NewLen(4)
	utils.GetSeesionMgr(c).Set("captcha", code)
	captcha.WriteImage(c.Writer, code, 100, 40)
	return
}

//验证码
func (this *AdminControl) logout(c *gin.Context) {
	utils.GetSeesionMgr(c).Set("sysInfo", nil)
	c.Redirect(200, "/")
	return
}
