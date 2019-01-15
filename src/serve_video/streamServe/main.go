package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"serve_video/streamServe/controler"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/videos/:vid-id", controler.StreamHandler)
	router.POST("/upload/:vid-id", controler.UploadHandler)
	router.GET("/testpage", controler.TestPageHandler)
	return router
}
func main() {
	r := RegisterHandlers()
	mh := controler.NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", mh)
}
