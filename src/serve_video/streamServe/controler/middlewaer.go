package controler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MiddleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := MiddleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func (m MiddleWareHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !m.l.GetConn() {
		SendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}
	m.r.ServeHTTP(w, req)
	defer m.l.RealseConn()
}
