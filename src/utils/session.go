package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
	"time"
)

var SessionId string

//Session数据
type Session struct {
	mLock    sync.RWMutex //互斥(保证线程安全)
	mExpire  int64        //垃圾回收时间
	mSession map[string]interface{}
}

//Session会话管理
type SessionMgr struct {
	mLock     sync.RWMutex //互斥(保证线程安全)
	mExpire  int64        //垃圾回收时间
	mSessions map[string]*Session
}

var sessionMgr *SessionMgr

//获取session管理器
func GetSeesionMgr(c *gin.Context) *SessionMgr {
	if sessionMgr == nil {
		sessionMgr = &SessionMgr{
			mSessions: make(map[string]*Session),
		}
	}
	var err error
	SessionId, err = c.Cookie("0376.kim")
	if err != nil || SessionId == "" {
		SessionId = Md5(strconv.Itoa(int(time.Now().UnixNano())))
	}
	//设置浏览器cookie
	c.SetCookie("0376.kim", SessionId, 300, "/", c.GetHeader("Host"), false, true)
	//垃圾回收
	go sessionMgr.gc()
	return sessionMgr
}

//设置session
func (this *SessionMgr) Set (name string,  value interface{}) {
	if s, ok := this.mSessions[name]; ok {
		s.mLock.RLock()
		defer s.mLock.RUnlock()
		s.mSession[name] = value
	}
}

//读取session
func (this *SessionMgr) Get(name string) interface{} {
	if s, ok := this.mSessions[SessionId]; ok {
		if v, ok := s.mSession[name]; ok {
			return v
		}
	}
	return nil
}

//GC回收
func (this *SessionMgr) gc() {
	this.mLock.Lock()
	defer this.mLock.Unlock()
	//定时清理
	if time.Now().Unix() > this.mExpire {
		for k, s := range this.mSessions {
			//删除超过时限的session
			if time.Now().Unix() > s.mExpire {
				delete(this.mSessions, k)
			}
		}
	}
	this.mExpire = time.Now().Unix()
	//定时回收
	time.AfterFunc(time.Duration(10) * time.Minute, func() {
		this.gc()
	})
}

