package controler

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

func CreateUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	io.WriteString(w, "hello, world")
}

func Login(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	usname := p.ByName("user_name")
	io.WriteString(w, usname)
}
