package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"serve_video/web/controler"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", controler.HomeHandler)
	//router.POST("/", homeHandler)
	//router.GET("/userhome", userHomeHandler)
	//router.POST("/userhome", userHomeHandler)
	//router.POST("/api", apiHandler)
	router.ServeFiles("/static/*filepath", http.Dir("./templates"))
	return router
}
func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}
