package session

import (
	"log"
	"serve_video/api/dbops"
	"serve_video/api/model"
	"serve_video/api/utils"
	"sync"
	"time"
)

//并发安全的Map机制，充当一个cache缓存
var sessionMap *sync.Map

func init() {
	// 并发安全Map，当作缓存 key-sessionId value-*Session实体
	sessionMap = &sync.Map{}
}

// 从数据库加载sessionId [无法理解为什么重复写入syncMap]
func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}
	r.Range(func(key, value interface{}) bool {
		ss := value.(*model.SimpleSession)
		sessionMap.Store(key, ss)
		return true
	})
}

// 产生新的sessionId 更新数据库和cache 返回session的id
func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := time.Now().Unix() // 返回时间戳 精确到秒
	ttl := ct + 30*60       //过期时间30分钟
	ss := &model.SimpleSession{Username: un, TTL: ttl}
	sessionMap.Store(id, ss)
	err := dbops.InsertSession(id, ttl, un)
	if err != nil {
		log.Print("存入session失败")
	}
	return id
}

// 判断sessionId是否过期， 过期就返回空和false，没过期就返回username和true
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := time.Now().Unix()
		// 如果session超时了
		if ss.(*model.SimpleSession).TTL < ct {
			DeleteExpiredSession(sid)
			return "", true
		}
		return ss.(*model.SimpleSession).Username, false
	}
	return "", true
}

// 从cache和数据库中删除session
func DeleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	err := dbops.DeleteSession(sid)
	if err != nil {
		log.Print("删除session失败")
	}
}
