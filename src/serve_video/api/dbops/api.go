package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"serve_video/api/model"
	"serve_video/api/utils"
	"time"
)

// 添加用户
func AddUserCredential(loginName string, pwd string) error {
	stmIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmIns.Close()
	return nil
}

// 返回密码
func GetUserCreadential(loginName string) (string, error) {
	stmOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}
	var pwd string
	err = stmOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmOut.Close()
	return pwd, nil
}

// 删除用户
func DeleteUserCreadential(loginName, pwd string) error {
	stmDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name = ? AND pwd = ?")
	if err != nil {
		log.Printf("DleteUser err :%s", err)
		return err
	}
	_, err = stmDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmDel.Close()
	return nil
}

// 添加video
func AddNewViedo(aid int, name string) (*model.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	// 固定格式 M D Y, HH:MM:SS [Format数值不能变]
	ctime := t.Format("Jan 02 2006, 15:04:05")
	stmIns, err := dbConn.Prepare("INSERT INTO video_info(id,author_id,name,display_ctime) VALUES (?,?,?,?)")
	if err != nil {
		return nil, err
	}
	_, err = stmIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}
	res := &model.VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtiem: ctime,
	}
	defer stmIns.Close()
	return res, nil
}

// 获取video
func GetVideoInfo(vid string) (*model.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")
	var aid int
	var dct string
	var name string
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer stmtOut.Close()
	res := &model.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtiem: dct}
	return res, nil
}

// 删除video
func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}
