package router

import (
	"control"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

var adminRouter *gin.Engine

//// 中间件
//func adminMiddleware(c *gin.Context) {
//	control.Session = utils.GetSeesionMgr(c)
//	c.Next()
//}

// 获取后台路由实例
func GetAdminRouter() *gin.Engine {
	// view目录地址
	viewPath, _ := filepath.Abs("./views/admin/*")
	c := control.GetAdminControl()

	r := gin.Default()
	// 加载views
	r.LoadHTMLGlob(viewPath)

	// 使用中间件、过滤session登录
	//r.Use(adminMiddleware)

	// 静态路由
	r.Static("/upload", "./upload") //上传目录
	r.Static("/static", "./static") //静态文件
	r.Static("/lib", "./lib") //静态文件

	// 基本路由
	r.GET("/index.html", c.Index)   //首页
	r.GET("/", c.Index)   //首页
	r.GET("/index", c.Index)   //首页
	r.GET("/welcome.html", c.Welcome) //欢迎页
	r.GET("/login.html", c.Login) // 登录页
	r.POST("/login.html", c.Login) // ajax登录
	r.GET("/captcha.html", c.Captcha) //验证码

	// 后台菜单
	r.GET("/sysmenu-list.html", c.SysMenuList) //后台菜单列表
	r.POST("/sysmenu-list.html", c.SysMenuList) //ajax获取后台菜单列表
	r.POST("/sysmenu-child.html", c.SysMenuChildList) //ajax后台菜单列表
	r.POST("/sysmenu-del.html", c.SysMenuDel) //删除后台菜单
	r.GET("/sysmenu-modify.html", c.SysMenuModify) //编辑后台菜单页
	r.POST("/sysmenu-modify.html", c.SysMenuModify) //ajax提交
	// 后台管理员角色
	r.GET("/sysrole-list.html", c.SysRoleList) //管理员角色列表
	r.POST("/sysrole-list.html", c.SysRoleList) //ajax获取管理员角色列表
	r.GET("/sysrole-modify.html", c.SysRoleModify) //编辑角色列表页
	r.POST("/sysrole-modify.html", c.SysRoleModify) //ajax提交

	// 后台管理员角色
	r.GET("/sysuser-list.html", c.SysUserList) //管理员列表
	r.POST("/sysuser-list.html", c.SysUserList) //ajax获取管理员列表
	r.POST("/sysuser-del.html", c.SysUserDel) //冻结、解冻管理员
	r.GET("/sysuser-modify.html", c.SysUserModify) //编辑管理员列表页
	r.POST("/sysuser-modify.html", c.SysUserModify) //ajax提交

	// 基础配置
	r.GET("/config-info.html", c.ConfigInfo) //基础配置页面
	r.POST("/config-info.html", c.ConfigInfo) //ajax获取基础配置详情

	return r
}
