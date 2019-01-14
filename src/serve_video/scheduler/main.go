package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"serve_video/scheduler/controller"
	"serve_video/scheduler/taskRunner"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/video-delete-record/:vid-id", controller.VideoDeleteHandler)
	return router
}
func main() {
	go taskRunner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}
