package controler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"serve_video/api/model"
	"serve_video/api/session"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UNAME = "X-User-Name"

// httpMiddleware 做身份验证之类的判断
type MiddleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := MiddleWareHandler{}
	m.r = r
	return &m
}

// 劫持http
func (m *MiddleWareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// check session
	validateUserSession(req)
	m.r.ServeHTTP(w, req)
}

// session检验
func validateUserSession(req *http.Request) bool {
	sid := req.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}
	usName, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	req.Header.Add(HEADER_FIELD_UNAME, usName)
	return true
}

// user检验
func validateUser(w http.ResponseWriter, req *http.Request) bool {
	usName := req.Header.Get(HEADER_FIELD_UNAME)
	if len(usName) == 0 {
		sendErrorRespinse(w, model.ErrorNotAuthUser)
		return false
	}

	return true
}
