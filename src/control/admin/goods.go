package admin

import (
	"github.com/gin-gonic/gin"
	"logic"
	"strconv"
	"strings"
)

// 商品分类列表
func (this *AdminControl) GoodsCategoryList (c *gin.Context) {
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
func (this *AdminControl) GoodsCategoryChildList (c *gin.Context) {
	if strings.ToUpper(c.Request.Method) == "POST" {
		parentId, _ := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
		if parentId == 0 {
			returnJson(c,1001,nil)
			return
		}
		iRelust := logic.GetGoodsCategory().GetGoodsCategoryChildListByParentId(int64(parentId))
		returnJson(c,0, iRelust)
		return
	}
	returnJson(c,1000,nil)
	return
}

// 删除商品分类
func (this *AdminControl) GoodsCategoryDel (c *gin.Context) {
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
func (this *AdminControl) GoodsCategoryModify (c *gin.Context) {
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
	// 获取商品分类列表
	if categoryId > 0 {
		info = logic.GetGoodsCategory().GetGoodsCategoryInfoById(int64(categoryId))
	}
	goodsCategoryTree := logic.GetGoodsCategory().GetGoodsCategoryListByLevel(1)
	returnHtml(c, "goods-category-modify.html", gin.H{
		"info": info,
		"goodsCategoryTree": goodsCategoryTree,
	});
	return
}