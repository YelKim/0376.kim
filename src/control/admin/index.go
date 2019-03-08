package admin

import (
	"conf"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"logic"
	"os"
	"runtime"
	"strings"
)

//首页
func (this *AdminControl) Index(c *gin.Context) {
	// 获取后台左侧菜单
	menuTree := logic.GetSysMenu().GetSysMenuListByLevel(2)
	returnHtml(c, "index.html", gin.H{"menuTree": menuTree})
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
		return
	}
	returnHtml(c, "login.html", nil)
	return
}

//验证码
func (this *AdminControl) Captcha(c *gin.Context) {
	code := captcha.NewLen(4)
	captcha.WriteImage(c.Writer, code, 100, 40)
	return
}
