package dbops

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var (
	dbConn *sql.DB
	err    error
)

// 复用dbConn，包含此包就会调用此初始化函数
func init() {
	dbConn, err = sql.Open("mysql",
		"root:ainiyu@/videoserve?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

}
