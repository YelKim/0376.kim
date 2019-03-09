package logic

import (
	"encoding/json"
	"html"
	"strings"
	"utils"
)

type Goods struct {
	Id           int32    `json:"id" bson:"id"`
	Name         string   `json:"name" bson:"name"`                   //商品名称
	Title        string   `json:"title" bson:"title"`                 //html title
	Keyword      string   `json:"keyword" bson:"keyword"`             //html keyword
	Description  string   `json:"description" bson:"description"`     //html description
	CostPrice    float64  `json:"cost_price" bson:"cost_price"`       //原价
	Stock        int64    `json:"stock" bson:"stock"`                 //库存
	Price        float64  `json:"price" bson:"price"`                 //库存
	Mainurl      string   `json:"mainurl" bson:"mainurl"`             //主图地址
	Imgurl       string   `json:"imgurl" bson:"imgurl"`               //图片地址用”,”隔开
	ImgurlArr    []string `json:"imgurl_arr" bson:"imgurl_arr"`       //图片地址列表
	Details      string   `json:"details" bson:"details"`             //商品详情
	TopLevel     int64    `json:"top_level" bson:"top_level"`         //商品顶级分类ID
	CategoryId   int64    `json:"category_id" bson:"category_id"`     //分类ID
	IsPlan       int64    `json:"is_plan" bson:"is_plan"`             //是否计划发布
	PlanAt       int64    `json:"plan_at" bson:"plan_at"`             //计划开始时间
	EndAt        int64    `json:"end_at" bson:"end_at"`               //计划结束时间
	Deleted      int      `json:"deleted" bson:"deleted"`             // 状态 0:正常 1: 冻结
	CategoryName string   `json:"category_name" bson:"category_name"` //分类名称
	UpdateAt     string   `json:"update_at" bson:"update_at"`
}

type goodsList struct {
	List  []*Goods `json:"list" bson:"list"`
	Total int64    `json:"total" bson:"total"`
}

// 分页获取商品列表
func (this *Goods) GetGoodListByPage (deleted, categoryId, page int64, keyword string) interface{} {
	jsonStr, _ := db.Call("Proc_Goods_pagination_v1.0", deleted, page, pageSize, keyword, categoryId)
	info := &goodsList{}
	json.Unmarshal([]byte(jsonStr), &info)
	iResult := make(map[string]interface{}, 0)
	iResult["total"] = info.Total
	iResult["pageSize"] = pageSize
	iResult["data"] = info.List
	return iResult
}

// 上架、下架商品
func (this *Goods) DelGoodsById (goodsId, deleted int64) int {
	jsonStr, _ := db.Call("Proc_Goods_delById_v1.0", goodsId, deleted)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询商品详情
func (this *Goods) GetGoodsInfoById (goodsId int64) *Goods {
	jsonStr, _ := db.Call("Proc_Goods_infoById_v1.0", goodsId)
	info := []*Goods{}
	json.Unmarshal([]byte(jsonStr), &info)
	info[0].ImgurlArr = strings.Split(info[0].Imgurl, ",")
	info[0].Details = html.UnescapeString(info[0].Details)
	return info[0]
}

// 添加、编辑商品
func (this *Goods) ModifyGoods (nickname, phone, password, remark string, roleId, sex, adminId int64) int {
	if len(password) > 0  {
		password = utils.Md5(password)
	}
	jsonStr, _ := db.Call("Proc_Goods_modify_v1.0", nickname, phone, sex, password, remark, roleId, adminId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}