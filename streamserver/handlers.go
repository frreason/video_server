package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, err := template.ParseFiles("./videos/upload.html")
	if err != nil {
		log.Printf("ParseFile error: %s", err)
	}

	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("streamHandler!")
	targetUrl := "http://frreason-videos.oss-cn-shenzhen.aliyuncs.com/videos/" + p.ByName("vid-id")
	http.Redirect(w, r, targetUrl, 301)

}

// func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	vid := p.ByName("vid-id")
// 	//vl := VIDEO_DIR + vid + ".mp4"
// 	vl := VIDEO_DIR + vid

// 	log.Println(vl)
// 	video, err := os.Open(vl)
// 	if err != nil {
// 		//log.Println("open error")
// 		log.Printf("open file error ")
// 		sendErrorResponse(w, http.StatusInternalServerError, "<h1> Internal error</hl>")
// 		return
// 	}

// 	w.Header().Set("Content-Type", "video/mp4")
// 	http.ServeContent(w, r, "", time.Now(), video) //以二进制流传输，客户端会以mp4格式解析

// 	defer video.Close()

// }

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error when try to get file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error ")
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %s", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error ")
	}
	fn := p.ByName("vid-id")                         //上传的文件名字
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666) //写入本地计算机，成功后再写入oss对象存储
	if err != nil {
		log.Printf("Write file error: %s", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal error ")
		return
	}

	ossfn := "videos/" + fn
	path := "./videos/" + fn
	bn := "frreason-videos"
	ret := UploadToOss(ossfn, path, bn)
	if !ret {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	os.Remove(path) //移除临时文件夹

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload successfully ")

}
