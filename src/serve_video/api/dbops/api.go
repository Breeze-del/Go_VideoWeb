package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
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
