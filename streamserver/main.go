package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandle(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/testpage", testPageHandler)
	router.GET("/videos/:vid-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)

	return router
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		log.Println("GetConn error")
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many request!")
		return
	}
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func main() {
	r := RegisterHandlers()

	mh := NewMiddleWareHandle(r, 2)

	http.ListenAndServe(":9000", mh)
}
