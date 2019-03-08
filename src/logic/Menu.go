package logic

import (
	"encoding/json"
	"sort"
)

type MenuInfo struct {
	Id          int32  `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`               //栏目名称
	Title       string `json:"title" bson:"title"`             //html title
	Keyword     string `json:"keyword" bson:"keyword"`         //seo关键字
	Description string `json:"description" bson:"description"` //seo描述
	Url         string `json:"url" bson:"url"`                 //跳转链接
	Sort        int32  `json:"sort" bson:"sort"`               //排序
	ParentId    int32  `json:"parent_id" bson:"parent_id"`     //无限分类 父类ID
	Level       int8   `json:"level" bson:"level"`             //层级
}

type Menu struct {
	MenuInfo
	UpateAt    int64 `json:"update_at" bson:"update_at"`
	ChildCount int64 `json:"child_count" bson:"child_count"`
}

type menuList struct {
	List  []*Menu `json:"list" bson:"list"`
	Total int64   `json:"total" bson:"total"`
}

// 分页获取栏目列表
func (this *Menu) GetMenuListByPage(page int32, keyword string) interface{} {
	jsonStr, _ := db.Call("Proc_Menu_pagination_v1.0", page, pageSize, keyword)
	info := &menuList{}
	json.Unmarshal([]byte(jsonStr), &info)
	iResult := make(map[string]interface{}, 0)
	iResult["total"] = info.Total
	iResult["pageSize"] = pageSize
	iResult["data"] = info.List
	return iResult
}

// 分页获取子栏目列表
func (this *Menu) GetMenuChildListByParentId(parentId int64) interface{} {
	jsonStr, _ := db.Call("Proc_Menu_child_v1.0", parentId)
	info := []*Menu{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info
}

// 删除栏目
func (this *Menu) DelMenuById(menuId int64) int {
	jsonStr, _ := db.Call("Proc_Menu_delById_v1.0", menuId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询栏目详情
func (this *Menu) GetMenuInfoById(menuId int64) *Menu {
	jsonStr, _ := db.Call("Proc_Menu_infoById_v1.0", menuId)
	info := []*Menu{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]
}

// 添加、编辑栏目
func (this *Menu) ModifyMenu(name, title, keyword, description, url string, parentId, sort, menuId int64) int {
	jsonStr, _ := db.Call("Proc_Menu_modify_v1.0", name, title, keyword, description, url, parentId, sort, menuId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

type menuTree struct {
	MenuInfo
	Children []*menuTree `json:"children" bson:"children"`
}

// 根据等级获取栏目列表
func (this *Menu) GetMenuListByLevel(level int) []*menuTree {
	jsonStr, _ := db.Call("Proc_Menu_listBylevel_v1.0", level)
	list := []*menuTree{}
	json.Unmarshal([]byte(jsonStr), &list)
	// 循环遍历list
	tempTree := make(map[int]map[int]*menuTree)
	for _, v := range list {
		if _, ok := tempTree[int(v.ParentId)]; !ok {
			tempTree[int(v.ParentId)] = make(map[int]*menuTree)
		}
		tempTree[int(v.ParentId)][int(v.Id)] = v
	}

	// 递归遍历
	var makeTreeCore func(root int, tree map[int]map[int]*menuTree) []*menuTree
	makeTreeCore = func(root int, tree map[int]map[int]*menuTree) []*menuTree {
		tmp := make(map[int32]*menuTree, 0)
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
		temp := make([]*menuTree, 0)
		for i := 0; i < len(keys); i++ {
			temp = append(temp, tmp[int32(keys[i])])
		}
		return temp
	}
	tree := makeTreeCore(0, tempTree)
	return tree
}
