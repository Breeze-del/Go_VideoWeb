package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"serve_video/api/controler"
)

// 返回一个router路由实例
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.POST("/user", controler.CreateUser)
	router.POST("/user/:user_name", controler.Login)
	return router
}
func main() {
	r := RegisterHandlers()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		os.Exit(2)
	}
}
