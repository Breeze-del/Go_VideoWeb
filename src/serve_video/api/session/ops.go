package session

import "sync"

//并发安全的Map机制，充当一个cache缓存
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

// 从数据库加载sessionId
func LoadSessionsFromDB() {

}

// 产生新的sessionId
func GenerateNewSessionId(un string) string {

}

// 判断sessionId是否过期， 过期就返回空和false，没过期就返回username和true
func IsSessionExpired(sid string) (string, bool) {

}
