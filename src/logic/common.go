package logic

import (
	"conf"
	"os"
	"path/filepath"
	"strings"
	"utils/mysql"
)

var db = mysql.GetMysql()

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