package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type IGoodsCategory interface {
	GetGoodsCategoryListByPage(int32, string) interface{}
	GetGoodsCategoryChildListByParentId(int64) interface{}
	DelGoodsCategoryById(int64) int
	GetGoodsCategoryInfoById(int64) *GoodsCategory
	ModifyGoodsCategory(string, string, int64, int64, int64) int
	GetGoodsCategoryListByLevel(int) []*goodsCategoryTree
}

type GoodsCategoryInfo struct {
	Id       int32
	Title    string                                    // 分类名称
	Imgurl   string                                    // 图标
	Sort     int32                                     //排序
	ParentId int32 `json:"parent_id" bson:"parent_id"` //无限分类 父类ID
	Level    int8                                      //层级
}

type GoodsCategory struct {
	GoodsCategoryInfo
	UpateAt    int64 `json:"update_at" bson:"update_at"`
	ChildCount int64 `json:"child_count" bson:"child_count"`
}

type goodsCategoryList struct {
	List  []*GoodsCategory
	Total int64
}

type goodsCategoryTree struct {
	GoodsCategoryInfo
	Children []*goodsCategoryTree
}

// 分页获取栏目列表
func (gc *GoodsCategory) GetGoodsCategoryListByPage(page int32, keyword string) interface{} {
	jsonStr, _ := db.Call("Proc_GoodsCategory_pagination_v1.0", page, pageSize, keyword)
	info := &goodsCategoryList{}
	json.Unmarshal([]byte(jsonStr), &info)
	iResult := make(map[string]interface{}, 0)
	iResult["total"] = info.Total
	iResult["pageSize"] = pageSize
	iResult["data"] = info.List
	return iResult
}

// 分页获取子栏目列表
func (gc *GoodsCategory) GetGoodsCategoryChildListByParentId(parentId int64) interface{} {
	jsonStr, _ := db.Call("Proc_GoodsCategory_child_v1.0", parentId)
	info := []*GoodsCategory{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info
}

// 删除栏目
func (gc *GoodsCategory) DelGoodsCategoryById(categotyId int64) int {
	jsonStr, _ := db.Call("Proc_GoodsCategory_delById_v1.0", categotyId)
	info := []map[string]string{}
	json.Unmarshal([]byte(jsonStr), &info)
	_type, _ := strconv.Atoi(info[0]["type"])
	if _type == 0 {
		// 判断文件是否存在
		if _, err := os.Stat("." + info[0]["imgurl"]); os.IsNotExist(err) == false {
			os.Remove("." + info[0]["imgurl"])
		}
	}
	return _type
}

// 根据ID查询后台菜单详情
func (gc *GoodsCategory) GetGoodsCategoryInfoById(categoryId int64) *GoodsCategory {
	jsonStr, _ := db.Call("Proc_GoodsCategory_infoById_v1.0", categoryId)
	info := []*GoodsCategory{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]
}

// 添加、编辑后台菜单
func (gc *GoodsCategory) ModifyGoodsCategory(title, imgurl string, parentId, sort, categoryId int64) int {
	// 判断是否有图标上传
	_imgurl := ""
	if len(imgurl) > 0 {
		_imgurl = moveUpfile(imgurl, "category")
	}
	jsonStr, _ := db.Call("Proc_GoodsCategory_modify_v1.0", title, parentId, _imgurl, sort, categoryId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	fmt.Println(info)
	return info[0]["type"]
}

// 根据等级获取菜单列表
func (gc *GoodsCategory) GetGoodsCategoryListByLevel(level int) []*goodsCategoryTree {
	jsonStr, _ := db.Call("Proc_GoodsCategory_listBylevel_v1.0", level)
	list := []*goodsCategoryTree{}
	json.Unmarshal([]byte(jsonStr), &list)
	// 循环遍历list
	tempTree := make(map[int]map[int]*goodsCategoryTree)
	for _, v := range list {
		if _, ok := tempTree[int(v.ParentId)]; !ok {
			tempTree[int(v.ParentId)] = make(map[int]*goodsCategoryTree)
		}
		tempTree[int(v.ParentId)][int(v.Id)] = v
	}

	// 递归遍历
	var makeTreeCore func(root int, tree map[int]map[int]*goodsCategoryTree) []*goodsCategoryTree
	makeTreeCore = func(root int, tree map[int]map[int]*goodsCategoryTree) []*goodsCategoryTree {
		tmp := make(map[int32]*goodsCategoryTree, 0)
		keys := []int{}
		i := 1
		for k, v := range tree[root] {
			if tree[k] != nil {
				v.Children = makeTreeCore(k, tree)
			}
			// 相同位置 自动+1
			if _, ok := tmp[v.Sort]; ok {
				v.Sort = v.Sort + int32(i)
			}
			tmp[v.Sort] = v
			keys = append(keys, int(v.Sort))
			i++
		}
		// 排序
		sort.Ints(keys)
		temp := make([]*goodsCategoryTree, 0)
		for i := 0; i < len(keys); i++ {
			temp = append(temp, tmp[int32(keys[i])])
		}
		return temp
	}
	tree := makeTreeCore(0, tempTree)
	return tree
}
