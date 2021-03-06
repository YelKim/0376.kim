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
func (ac *AdminControl) Index(c *gin.Context) {
	// 获取后台左侧菜单
	menuTree := logic.ISysMenu(sysmenu).GetSysMenuListByLevel(2)
	sessionId, _ := c.Cookie("session_id")
	sysInfo := utils.GetSeesionMgr(c).Get(sessionId,"sysInfo")
	returnHtml(c, "index.html", gin.H{
		"menuTree": menuTree,
		"info": sysInfo,
	})
	return
}

//欢迎页
func (ac *AdminControl) Welcome(c *gin.Context) {
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
func (ac *AdminControl) Login(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		code := c.DefaultPostForm("captcha", "")
		account := c.DefaultPostForm("account", "")
		password := c.DefaultPostForm("password", "")
		sessionId, _ := c.Cookie("session_id")
		reald := utils.GetSeesionMgr(c).Get(sessionId,"captcha")
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
		iResult := logic.ISysUser(sysuser).Login(account, utils.Md5(strings.Trim(password, "")), c)
		returnJson(c, iResult, nil)
		return
	}
	returnHtml(c, "login.html", nil)
	return
}

//验证码
func (ac *AdminControl) Captcha(c *gin.Context) {
	code := captcha.NewLen(4)
	sessionId, _ := c.Cookie("session_id")
	utils.GetSeesionMgr(c).Set(sessionId, "captcha", code)
	captcha.WriteImage(c.Writer, code, 100, 40)
	return
}

//验证码
func (ac *AdminControl) Logout(c *gin.Context) {
	sessionId, _ := c.Cookie("session_id")
	utils.GetSeesionMgr(c).Set(sessionId, "sysInfo", nil)
	c.Redirect(200, "/")
	return
}
