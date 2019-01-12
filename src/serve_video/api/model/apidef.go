package model

// request
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

// response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

// Data model
type VideoInfo struct {
	Id           string
	AuthorId     int
	Name         string
	DisplayCtiem string
}

// 评论实体
type Comments struct {
	Id      string
	VideoId string
	Author  string
	Content string
}

// 从db里拿sessionId 登陆分配sessionId 检验session是否过期
type SimpleSessiong struct {
	UserName string
	TTL      int64
}
