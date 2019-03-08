package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"reflect"
)

type Config struct {
	Id      int32  `json:"id" bson:"id"`
	Type    int32  `json:"type" bson:"type"`       //类型  1：基础配置
	Content string `json:"content" bson:"content"` // 配置内容json字符串
	UpateAt int64  `json:"update_at" bson:"update_at"`
}

// 全局SEO配置
type SeoConfig struct {
	Keyword    string `json:"keyword" bson:"keyword"`          // 关键字
	Description string `json:"description" bson:"description"` // 描述
}

// 全局SEO配置
type InfoConfig struct {
	Platform  string `json:"platform" bson:"platform"`   // 关键字
	RecordNo string `json:"record_no" bson:"record_no"` // 版权
	Wx        string `json:"wx" bson:"wx"`               // 公众号
	Tel       string `json:"tel" bson:"tel"`             // 电话
	Address   string `json:"address" bson:"address"`     // 公司地址
}


// 添加、编辑基础配置
func (this *Config) ModifyConfig (_type int, content string) int {
	jsonStr, _ := db.Call("Proc_Config_modify_v1.0", _type, content)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询管理员角色详情
func (this *Config) GetConfigInfoByType (_type int) map[string]interface{} {
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

