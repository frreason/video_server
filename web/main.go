package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)
	router.POST("/", homeHandler)
	router.GET("/help", helpHandler)
	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)
	router.GET("/videos", videosHandler)
	router.GET("/help", myHelpHandler)
	router.POST("/api", apiHandler)
	router.POST("/upload/:vid-id", proxyUploadHandler)
	router.GET("/videos/:vid-id", proxyVideoHandler)
	router.ServeFiles("/statics/*filepath", http.Dir("./template"))
	return router
}

func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}
