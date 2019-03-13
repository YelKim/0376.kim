package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"reflect"
)

type IConfig interface {
	ModifyConfig(int, string) int
	GetConfigInfoByType(int) map[string]interface{}
}

type Config struct {
	Id      int32
	Type    int32  //类型  1：基础配置
	Content string // 配置内容json字符串
	UpateAt int64 `json:"update_at" bson:"update_at"`
}

// 全局SEO配置
type SeoConfig struct {
	Keyword     string // 关键字
	Description string // 描述
}

// 全局SEO配置
type InfoConfig struct {
	Platform string                                     // 关键字
	RecordNo string `json:"record_no" bson:"record_no"` // 版权
	Wx       string                                     // 公众号
	Tel      string                                     // 电话
	Address  string                                     // 公司地址
}


// 添加、编辑基础配置
func (c *Config) ModifyConfig (_type int, content string) int {
	jsonStr, _ := db.Call("Proc_Config_modify_v1.0", _type, content)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询管理员角色详情
func (c *Config) GetConfigInfoByType (_type int) map[string]interface{} {
	jsonStr, _ := db.Call("Proc_Config_info_v1.0", _type)
	info := []*Config{}
	json.Unmarshal([]byte(jsonStr), &info)
	var content interface{}

	switch _type {
	case 1:
		content = &InfoConfig{}
		break
	case 2:
		content = &SeoConfig{}
	}
	json.Unmarshal([]byte(info[0].Content), &content)
	// 获取结构体 json tag
	var keys []string
	t := reflect.TypeOf(content)
	for i := 0; i < t.Elem().NumField(); i++ {
		keys = append(keys, t.Elem().Field(i).Tag.Get("json"))
	}
	return gin.H{
		"content": content,
		"keys": keys,
	}
}

