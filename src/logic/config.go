package logic

import (
	"encoding/json"
)

type Config struct {
	Id      int32  `json:"id" bson:"id"`
	Type    int32  `json:"type" bson:"type"`       //类型  1：基础配置
	Content string `json:"content" bson:"content"` // 配置内容json字符串
	UpateAt int64  `json:"update_at" bson:"update_at"`
}

// 添加、编辑基础配置
func (this *Config) ModifyConfig (_type int, content string) int {
	jsonStr, _ := db.Call("Proc_Config_modify_v1.0", _type, content)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}
