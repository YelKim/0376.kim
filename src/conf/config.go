package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
	"os"
	"path/filepath"
)

//redis 配置
type Redis struct {
	Conn   string
	Passwd string
}

type Config struct {
	Redis     Redis
	Mysql     string
	StaticDir string
	FileSize  int32
	ApiPort   string
	AdminPort string
	FileExt   []string
	PageSize  int
}

var opts *Config
var fileModTime int64

//获取配置信息
func GetConfig() *Config {
	err := ParseToml()
	if err != nil {
		glog.Fatalf("Profile error: %v\n", err)
	}
	return opts
}

//读取配置
func ParseToml() error {
	flag.Parse()
	filePath, _ := filepath.Abs("./config.toml")
	file, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return err
	}
	if opts == nil || file.ModTime().Unix() > fileModTime {
		glog.Infoln("Parse profile...")
		var conf Config
		if _, err := toml.DecodeFile(filePath, &conf); err != nil {
			return err
		}
		fileModTime = file.ModTime().Unix()
		opts = &conf
		glog.Infof("Config info: %v\n", err)
	}
	return nil
}