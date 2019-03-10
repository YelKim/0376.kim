package admin

import (
	"github.com/gin-gonic/gin"
	"logic"
	"strconv"
	"strings"
	"time"
)

// 商品分类列表
func (this *AdminControl) GoodsCategoryList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		keyword := c.DefaultPostForm("keyword", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))
		iRelust := logic.GetGoodsCategory().GetGoodsCategoryListByPage(int32(page), keyword)
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	returnHtml(c, "goods-category.html", gin.H{
		"keyword": keyword,
	})
	return
}

// 子商品分类
func (this *AdminControl) GoodsCategoryChildList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
		if parentId == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.GetGoodsCategory().GetGoodsCategoryChildListByParentId(int64(parentId))
		returnJson(c, 0, iRelust)
		return
	}
	returnJson(c, 1000, nil)
	return
}

// 删除商品分类
func (this *AdminControl) GoodsCategoryDel(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		categoryId, _ := strconv.Atoi(c.DefaultPostForm("category_id", "0"))
		if categoryId == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.GetGoodsCategory().DelGoodsCategoryById(int64(categoryId))
		returnJson(c, iRelust, nil)
		return
	}
	returnJson(c, 1000, nil)
	return
}

// 编辑商品分类
func (this *AdminControl) GoodsCategoryModify(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		categoryId, _ := strconv.Atoi(c.DefaultPostForm("category_id", "0"))
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
		sort, _ := strconv.Atoi(c.DefaultPostForm("sort", "999"))
		title := c.DefaultPostForm("title", "")
		imgurl := c.DefaultPostForm("imgurl", "")
		iResult := logic.GetGoodsCategory().ModifyGoodsCategory(title, imgurl, int64(parentId), int64(sort), int64(categoryId))
		returnJson(c, iResult, nil)
		return
	}
	// 获取商品分类详情
	categoryId, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	info := &logic.GoodsCategory{}
	info.Sort = 999
	// 获取商品分类列表
	if categoryId > 0 {
		info = logic.GetGoodsCategory().GetGoodsCategoryInfoById(int64(categoryId))
	}
	goodsCategoryTree := logic.GetGoodsCategory().GetGoodsCategoryListByLevel(1)
	returnHtml(c, "goods-category-modify.html", gin.H{
		"info":              info,
		"goodsCategoryTree": goodsCategoryTree,
	})
	return
}

// 商品列表
func (this *AdminControl) GoodsList(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		keyword := c.DefaultPostForm("keyword", "")
		page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))
		deleted, _ := strconv.Atoi(c.DefaultPostForm("deleted", "0"))
		categoryId, _ := strconv.Atoi(c.DefaultPostForm("category_id", "-1"))
		iRelust := logic.GetGoods().GetGoodListByPage(int64(deleted), int64(categoryId), int64(page), keyword)
		returnJson(c, 0, iRelust)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	deleted, _ := strconv.Atoi(c.DefaultQuery("deleted", "0"))
	categoryId, _ := strconv.Atoi(c.DefaultQuery("category_id", "-1"))
	//查询商品分类
	goodsCategoryTree := logic.GetGoodsCategory().GetGoodsCategoryListByLevel(2)
	returnHtml(c, "goods.html", gin.H{
		"goodsCategoryTree": goodsCategoryTree,
		"keyword": keyword,
		"deleted": deleted,
		"categoryId": categoryId,
	})
	return
}

// 上、下架商品
func (this *AdminControl) GoodsDel(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		goodsId, _ := strconv.Atoi(c.DefaultPostForm("goods_id", "0"))
		deleted, _ := strconv.Atoi(c.DefaultPostForm("deleted", "-1"))
		if goodsId == 0 || deleted == -1 {
			returnJson(c, 1001, nil)
			return
		}
		iRelust := logic.GetGoods().DelGoodsById(int64(goodsId), int64(deleted))
		returnJson(c, iRelust, nil)
		return
	}
	returnJson(c, 1000, nil)
	return
}


// 编辑商品
func (this *AdminControl) GoodsModify(c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		goodsId, _ := strconv.Atoi(c.DefaultPostForm("goods_id", "0"))
		categoryId, _ := strconv.Atoi(c.DefaultPostForm("category_id", "0"))
		stock, _ := strconv.Atoi(c.DefaultPostForm("stock", "0"))
		isPlan, _ := strconv.Atoi(c.DefaultPostForm("is_plan", "0"))
		adminId := 0
		name := c.DefaultPostForm("name", "")
		title := c.DefaultPostForm("title", "")
		keyword := c.DefaultPostForm("keyword", "")
		description := c.DefaultPostForm("description", "")
		cost_price,_ := strconv.ParseFloat(c.DefaultPostForm("cost_price", "0"), 64)
		price,_ := strconv.ParseFloat(c.DefaultPostForm("price", "0"), 64)
		details := c.DefaultPostForm("details", "")
		planAt := c.DefaultPostForm("plan_at", "")
		if len(planAt) > 0 {
			tm, _ := time.Parse("2006-01-02", planAt)
			planAt = strconv.Itoa(int(tm.Unix()))
		}
		endAt := c.DefaultPostForm("end_at", "2099-01-01")
		tm, _ := time.Parse("2006-01-02", endAt)
		endAt = strconv.Itoa(int(tm.Unix()))
		mainurl := c.DefaultPostForm("mainurl", "")
		imgurl := strings.Trim(c.DefaultPostForm("imgurl", ""), ",")
		if price == 0 ||  cost_price == 0 || len(mainurl) == 0 || categoryId == 0 || len(name) == 0 || len(details) == 0 {
			returnJson(c, 1001, nil)
			return
		}
		iResult := logic.GetGoods().ModifyGoods(name, title, keyword, description, cost_price, price, details, mainurl, imgurl, int64(categoryId), int64(stock), int64(adminId), int64(isPlan), planAt, endAt, int64(goodsId))
		returnJson(c, iResult, nil)
		return
	}
	// 获取商品详情
	goodsId, _ := strconv.Atoi(c.DefaultQuery("goods_id", "0"))
	info := &logic.Goods{}
	// 获取商品分类列表
	if goodsId > 0 {
		info = logic.GetGoods().GetGoodsInfoById(int64(goodsId))
	}
	goodsCategoryTree := logic.GetGoodsCategory().GetGoodsCategoryListByLevel(2)
	returnHtml(c, "goods-modify.html", gin.H{
		"info":              info,
		"goodsCategoryTree": goodsCategoryTree,
	})
	return
}