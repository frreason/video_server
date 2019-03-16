package dbops

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	//这里已经提前连接好了数据库，并且使用video_server数据库！！！
	dbConn, err = sql.Open("mysql", "root:061365404abc@tcp(localhost:3306)/video_server?charset=utf8")

	if err != nil {
		panic(err.Error())
	}
}
