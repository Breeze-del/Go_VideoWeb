package dbops

import (
	"database/sql"
	"log"
	"serve_video/api/model"
	"strconv"
	"sync"
)

// 向数据库读写session
func InsertSession(sid string, ttl int64, usname string) error {
	ttlStr := strconv.FormatInt(ttl, 10)
	// 预加载 防止sql注入攻击
	stmIns, err := dbConn.Prepare("INSERT INTO sessions(session_id,TTL,login_name) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmIns.Exec(sid, ttlStr, usname)
	if err != nil {
		return err
	}
	defer stmIns.Close()
	return nil
}

// 根据sessionId 查session表
func RetrieveSession(sid string) (*model.SimpleSessiong, error) {
	ss := &model.SimpleSessiong{}
	stmOut, err := dbConn.Prepare("SELECT sessions.login_name,sessions.TTL FROM sessions WHERE sessions.session_id=?")
	if err != nil {
		return nil, err
	}
	var usname string
	var ttl string
	err = stmOut.QueryRow(sid).Scan(&usname, &ttl)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stmOut.Close()
	ttlInt64, err := strconv.ParseInt(ttl, 10, 64)
	if err != nil {
		return nil, err
	}
	ss.UserName = usname
	ss.TTL = ttlInt64
	return ss, nil
}

// 获取所有session,并返回syncMap当作缓存
func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	for rows.Next() {
		var id string
		var ttlstr string
		var login_name string
		if err := rows.Scan(&id, &ttlstr, &login_name); err != nil {
			log.Printf("retrive sessions error: %s", err)
			break
		}
		if ttl, err1 := strconv.ParseInt(ttlstr, 10, 64); err1 == nil {
			ss := &model.SimpleSessiong{UserName: login_name, TTL: ttl}
			m.Store(id, ss)
			log.Printf(" session id: %s, ttl: %d", id, ss.TTL)
		}
	}
	defer stmtOut.Close()
	return m, nil
}

// 根据sessionId删除session
func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	if _, err := stmtOut.Exec(sid); err != nil {
		return err
	}
	stmtOut.Close()
	return nil
}
