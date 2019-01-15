package dbops

import (
	"database/sql"
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

func GetUser(loginName string) (*model.User, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE login_name=?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	var id int
	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	res := &model.User{Id: id, LoginName: loginName, Pwd: pwd}
	defer stmtOut.Close()
	return res, nil
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
		DisplayCtime: ctime,
	}
	defer stmIns.Close()
	return res, nil
}

// 获取video
func GetVideoInfo(vid string) (*model.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")
	res := model.VideoInfo{Id: vid}
	err = stmtOut.QueryRow(vid).Scan(&res.AuthorId, &res.Name, &res.DisplayCtime)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	stmtOut.Close()
	return &res, nil
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
	stmtDel.Close()
	return nil
}

// 添加评论
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}
	stmIns, err := dbConn.Prepare("INSERT INTO comments(id,video_id,author_id,content) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	stmIns.Close()
	return nil
}

// 根据视频id 找出某个时间内所有的评论用户姓名和评论
func ListComments(vid string, from, to int) ([]*model.Comment, error) {
	stmOut, err := dbConn.Prepare(`SELECT comments.id,users.login_name,comments.content from comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)`)
	var res []*model.Comment
	rows, err := stmOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		comment := model.Comment{}
		if err := rows.Scan(&comment.Id, &comment.Author, &comment.Content); err != nil {
			return res, err
		}
		res = append(res, &comment)
	}
	defer stmOut.Close()
	return res, nil
}

func ListVideoInfo(uname string, from, to int) ([]*model.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.display_ctime FROM video_info
		INNER JOIN users ON video_info.author_id = users.id
		WHERE users.login_name=? AND video_info.create_time > FROM_UNIXTIME(?) AND video_info.create_time<=FROM_UNIXTIME(?)`)
	var res []*model.VideoInfo
	if err != nil {
		return res, err
	}
	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		log.Printf("%s", err)
		return res, err
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}
		vi := &model.VideoInfo{Id: id, AuthorId: aid, Name: name, DisplayCtime: ctime}
		res = append(res, vi)
	}
	defer stmtOut.Close()
	return res, nil
}
