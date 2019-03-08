package logic

import (
	"conf"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"utils/mysql"
)

var db = mysql.GetMysql()
var once sync.Once
var pageSize = conf.GetConfig().PageSize

// 把文件从临时目录移动到新文件
func moveUpfile(fileName, folder string) string {
	//判断文件夹是否存在
	newFileName := ""
	if _, err := os.Stat("." + fileName); os.IsNotExist(err) {
		return newFileName
	}
	// 判断新文件夹都是否存在
	if strings.Index(fileName, "/tmp/") > 0 {
		newFileName = strings.Replace(fileName, "/tmp/", "/" + folder + "/", 1)
		filePath, _ := filepath.Split(newFileName)
		if _, err := os.Stat( "." + filePath); os.IsNotExist(err) {
			os.MkdirAll("." + filePath, 0755)
		}
		err := os.Rename("." + fileName, "." + newFileName)
		if nil != err {
			return ""
		}
	}
	return newFileName
}

// 后台菜单
var sysMenu *SysMenu
func GetSysMenu() *SysMenu {
	once.Do(func() {
		sysMenu = &SysMenu{}
	})
	return sysMenu
}

// 管理员角色
var sysRole *SysRole
func GetSysRole() *SysRole {
	once.Do(func() {
		sysRole = &SysRole{}
	})
	return sysRole
}

// 管理员
var sysUser *SysUser
func GetSysUser() *SysUser {
	once.Do(func() {
		sysUser = &SysUser{}
	})
	return sysUser
}

// 基础配置
var config *Config
func GetConfig() *Config {
	once.Do(func() {
		config = &Config{}
	})
	return config
}

// 栏目管理
var menu *Menu
func GetMenu() *Menu {
	once.Do(func() {
		menu = &Menu{}
	})
	return menu
}

// 商品管理
var goodsCategory *GoodsCategory
func GetGoodsCategory() *GoodsCategory {
	once.Do(func() {
		goodsCategory = &GoodsCategory{}
	})
	return goodsCategory
}


