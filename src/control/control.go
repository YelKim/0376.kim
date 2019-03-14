package control

import (
	"control/admin"
	"sync"
	"utils"
)

var once sync.Once
var adminControl *admin.AdminControl
var Session *utils.SessionMgr

func GetAdminControl() *admin.AdminControl {
	once.Do(func() {
		adminControl = &admin.AdminControl{}
	})
	return adminControl
}