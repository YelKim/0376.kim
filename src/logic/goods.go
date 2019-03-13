package logic

import (
	"encoding/json"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"
)

type IGoods interface {
	GetGoodListByPage(int64, int64, int64, string) interface{}
	DelGoodsById(int64, int64) int
	GetGoodsInfoById(int64) *Goods
	ModifyGoods(string, string, string, string, float64, float64, string, string, string, int64, int64, int64, int64, string, string, int64) int
}

type Goods struct {
	Id           int32
	Name         string                                             //商品名称
	Title        string                                             //html title
	Keyword      string                                             //html keyword
	Description  string                                             //html description
	Stock        int64                                              //库存
	CostPrice    string `json:"cost_price" bson:"cost_price"`       //原价
	Price        string                                             //价格
	Mainurl      string                                             //主图地址
	Imgurl       string                                             //图片地址用”,”隔开
	ImgurlArr    []string `json:"imgurl_arr" bson:"imgurl_arr"`     //图片地址列表
	ImgurlLen    int      `json:"imgurl_len" bson:"imgurl_len"`     //图片地址列表
	Details      string                                             //商品详情
	TopLevel     int64  `json:"top_level" bson:"top_level"`         //商品顶级分类ID
	CategoryId   int64  `json:"category_id" bson:"category_id"`     //分类ID
	IsPlan       int64  `json:"is_plan" bson:"is_plan"`             //是否计划发布
	PlanAt       int64  `json:"plan_at" bson:"plan_at"`             //计划开始时间
	PlanAtTxt    string `json:"plan_at_txt" bson:"plan_at_txt"`     //计划开始时间
	EndAt        int64  `json:"end_at" bson:"end_at"`               //计划结束时间
	EndAtTxt     string `json:"end_at_txt" bson:"end_at_txt"`       //计划结束时间始时间
	Deleted      int                                                // 状态 0:正常 1: 冻结
	CategoryName string `json:"category_name" bson:"category_name"` //分类名称
	UpdateAt     string `json:"update_at" bson:"update_at"`
}

type goodsList struct {
	List  []*Goods
	Total int64
}

// 分页获取商品列表
func (g *Goods) GetGoodListByPage (deleted, categoryId, page int64, keyword string) interface{} {
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
func (g *Goods) DelGoodsById (goodsId, deleted int64) int {
	jsonStr, _ := db.Call("Proc_Goods_delById_v1.0", goodsId, deleted)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询商品详情
func (g *Goods) GetGoodsInfoById (goodsId int64) *Goods {
	jsonStr, _ := db.Call("Proc_Goods_infoById_v1.0", goodsId)
	info := []*Goods{}
	json.Unmarshal([]byte(jsonStr), &info)
	info[0].ImgurlArr = append(info[0].ImgurlArr, info[0].Mainurl)
	imgArr := strings.Split(info[0].Imgurl, ",")
	for _, v := range imgArr {
		if len(v) > 0 {
			info[0].ImgurlArr = append(info[0].ImgurlArr, v)
		}
	}
	if info[0].PlanAt > 0 {
		info[0].PlanAtTxt = time.Unix(info[0].PlanAt, 0).Format("2006-01-02")
	}
	if info[0].EndAt > 0 {
		info[0].EndAtTxt = time.Unix(info[0].EndAt, 0).Format("2006-01-02")
	}
	info[0].ImgurlLen = len(info[0].ImgurlArr)
	info[0].Details = html.UnescapeString(info[0].Details)
	price, _ := strconv.ParseFloat(info[0].Price, 64)
	info[0].Price = fmt.Sprintf("%.2f", price)
	cost_price, _ := strconv.ParseFloat(info[0].CostPrice, 64)
	info[0].CostPrice = fmt.Sprintf("%.2f", cost_price)
	// 转换价格
	return info[0]
}

// 添加、编辑商品
func (g *Goods) ModifyGoods (name, title, keyword, description string, cost_price, price float64, details, mainImg, imgurl string, categoryId, stock, adminId, isPlan int64, planAt, endAt string, goodsId int64) int {
	//处理主图
	_mainImg := ""
	if len(mainImg) > 0 {
		_mainImg = moveUpfile(mainImg, "goods")
	}
	// 处理图片列表
	imgurlArr := strings.Split(imgurl, ",")
	_imgurlArr := []string{}
	for _, v := range imgurlArr {
		if strings.Trim(v, "") == mainImg {
			continue
		}
		if strings.Index(v, "/tmp/") > 0 {
			v = moveUpfile(v, "goods")
		}
		if len(v) > 0 {
			_imgurlArr = append(_imgurlArr, v)
		}
	}
	jsonStr, _ := db.Call("Proc_Goods_modify_v1.0", name, title, keyword, description, cost_price, price, categoryId, stock, details, _mainImg, strings.Join(_imgurlArr, ","), isPlan, planAt, endAt, goodsId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}