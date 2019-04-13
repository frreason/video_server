package main

import (
	"bytes"
	"encoding/json"
	"frreason/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

//将各种http请求分发到各个真正得端口得api
func request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	u, _ := url.Parse(b.Url)
	u.Host = config.GetLbAddr() + ":" + u.Port()
	newUrl := u.String()
	switch b.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("httpClient.Do error: %s \n", err)
			return
		}
		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", newUrl, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("httpClient.Do error: %s\n", err)
			log.Printf("     ")
			return
		}
		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("httpClient.Do error: %s \n", err)
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
		return
	}

}

func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		re, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(re))
		return
	}
	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}
