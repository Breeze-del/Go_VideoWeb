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
	// 注册路由
	router.POST("/user", controler.CreateUser)
	router.POST("/user/:username", controler.Login)
	router.GET("/user/:username", controler.GetUserInfo)
	router.POST("/user/:username/videos", controler.AddNewVideo)
	router.GET("/user/:username/videos", controler.ListAllVideos)
	router.DELETE("/user/:username/videos/:vid-id", controler.DeleteVideo)
	router.POST("/videos/:vid-id/comments", controler.PostComment)
	router.GET("/videos/:vid-id/comments", controler.ShowComments)
	return router
}
func main() {
	// 返回一个httpRouter
	r := RegisterHandlers()
	// 先将router放入结构体，然后结构体实现handler
	// 劫持掉httpResponseWriter和 Request
	mh := controler.NewMiddleWareHandler(r)
	err := http.ListenAndServe(":8000", mh)
	if err != nil {
		os.Exit(2)
	}
}
