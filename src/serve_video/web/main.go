package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"serve_video/web/controler"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", controler.HomeHandler)
	router.POST("/", controler.HomeHandler)
	router.GET("/userhome", controler.UserHomeHandler)
	router.POST("/userhome", controler.UserHomeHandler)
	router.POST("/api", controler.ApiHandler)
	router.ServeFiles("/static/*filepath", http.Dir("./templates"))
	return router
}
func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}
