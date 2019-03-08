package logic

import (
	"encoding/json"
	"html"
	"sort"
)

type SysMenuInfo struct {
	Id       int32         `json:"id" bson:"id"`
	Title    string        `json:"title" bson:"title"`         //菜单名称
	Icon     string `json:"icon" bson:"icon"`           //一级菜单图标
	Control  string        `json:"control" bson:"control"`     //controller
	Action   string        `json:"action" bson:"action"`       // action
	Sort     int32         `json:"sort" bson:"sort"`           //排序
	ParentId int32         `json:"parent_id" bson:"parent_id"` // 无限分类 父类ID
	Level    int8          `json:"level" bson:"level"`         // 层级
}

type SysMenu struct {
	SysMenuInfo
	UpateAt    int64 `json:"update_at" bson:"update_at"`
	ChildCount int64 `json:"child_count" bson:"child_count"`
}

type sysMenuList struct {
	List  []*SysMenu `json:"list" bson:"list"`
	Total int64      `json:"total" bson:"total"`
}

// 分页获取后台菜单列表
func (this *SysMenu) GetSysMenuListByPage (page int32, keyword  string) interface{} {
	jsonStr, _ := db.Call("Proc_SysMenu_pagination_v1.0", page, pageSize, keyword)
	info := &sysMenuList{}
	json.Unmarshal([]byte(jsonStr), &info)
	iResult := make(map[string]interface{}, 0)
	iResult["total"] = info.Total
	iResult["pageSize"] = pageSize
	iResult["data"] = info.List
	return iResult
}

// 分页获取后台子菜单列表
func (this *SysMenu) GetSysMenuChildListByParentId (parentId int64) interface{} {
	jsonStr, _ := db.Call("Proc_SysMenu_child_v1.0", parentId)
	info := []*SysMenu{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info
}

// 删除后台菜单
func (this *SysMenu) DelSysMenuById (menuId int64) int {
	jsonStr, _ := db.Call("Proc_SysMenu_delById_v1.0", menuId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

// 根据ID查询后台菜单详情
func (this *SysMenu) GetSysMenuInfoById (menuId int64) *SysMenu {
	jsonStr, _ := db.Call("Proc_SysMenu_infoById_v1.0", menuId)
	info := []*SysMenu{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]
}

// 添加、编辑后台菜单
func (this *SysMenu) ModifySysMenu (menuId, parentId, sort int64, title, icon, control, action string) int {
	jsonStr, _ := db.Call("Proc_SysMenu_modify_v1.0", title, icon, control, action, parentId, sort, menuId)
	info := []map[string]int{}
	json.Unmarshal([]byte(jsonStr), &info)
	return info[0]["type"]
}

type sysMenuTree struct {
	SysMenuInfo
	Children []*sysMenuTree `json:"children" bson:"children"`
}

// 根据等级获取菜单列表
func (this *SysMenu) GetSysMenuListByLevel(level int) []*sysMenuTree {
	jsonStr, _ := db.Call("Proc_SysMenu_listBylevel_v1.0", level)
	list := []*sysMenuTree{}
	json.Unmarshal([]byte(jsonStr), &list)
	// 循环遍历list
	tempTree := make(map[int]map[int]*sysMenuTree)
	for _, v := range list {
		if _, ok := tempTree[int(v.ParentId)]; !ok {
			tempTree[int(v.ParentId)] = make(map[int]*sysMenuTree)
		}
		tempTree[int(v.ParentId)][int(v.Id)] = v
	}

	// 递归遍历
	var makeTreeCore func(root int, tree map[int]map[int]*sysMenuTree) []*sysMenuTree
	makeTreeCore = func(root int, tree map[int]map[int]*sysMenuTree) []*sysMenuTree {
		tmp := make(map[int32]*sysMenuTree, 0)
		keys := []int{}
		i := 1
		for k, v := range tree[root] {
			if tree[k] != nil {
				v.Icon = html.UnescapeString(v.Icon)
				v.Children = makeTreeCore(k, tree)
			}
			// 相同位置 自动+1
			if _, ok := tmp[v.Sort]; ok  {
				v.Sort = v.Sort + int32(i)
			}
			tmp[v.Sort] = v
			keys = append(keys, int(v.Sort))
			i++;
		}
		// 排序
		sort.Ints(keys)
		temp := make([]*sysMenuTree, 0)
		for i := 0; i < len(keys); i++ {
			temp = append(temp, tmp[int32(keys[i])])
		}
		return temp
	}
	tree := makeTreeCore(0, tempTree)
	return tree
}