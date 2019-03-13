package utils

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//Session数据
type Session struct {
	mSessionId string
	mLastTime  int64
	mValues    map[string]interface{}
}

//Session会话管理
type SessionMgr struct {
	mLock     sync.RWMutex //互斥(保证线程安全)
	mExpire   int64        //定时器
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
	sessionId, err := c.Cookie("session_id")
	if err != nil || sessionId == "" {
		sessionId = Md5(strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(rand.Intn(10)))[8:24]
	}
	//设置浏览器cookie
	c.SetCookie("session_id", sessionId, 300, "/", c.GetHeader("Host"), false, true)
	//垃圾回收
	go sessionMgr.gc()
	return sessionMgr
}

//设置session
func (this *SessionMgr) Set (sessionId, name string, value interface{}) {
	this.mLock.RLock()
	defer this.mLock.RUnlock()
	if s, ok := this.mSessions[sessionId]; ok {
		s.mValues[name] = value
		s.mLastTime = time.Now().Unix()
	} else {
		val := make(map[string]interface{})
		val[name] = value
		session := &Session{
			mSessionId: sessionId,
			mLastTime:  time.Now().Unix(),
			mValues:    val,
		}
		this.mSessions[sessionId] = session
	}
}

//读取session
func (this *SessionMgr) Get(sessionId, name string) interface{} {
	if s, ok := this.mSessions[sessionId]; ok {
		if v, ok := s.mValues[name]; ok {
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
			if time.Now().Unix() - s.mLastTime == 1800 {
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

