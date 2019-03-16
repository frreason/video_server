package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/session"
	"video_server/api/utils"

	"github.com/julienschmidt/httprouter"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, SessionId: id}
	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}

// //对用户登陆的验证
// func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	res, _ := ioutil.ReadAll(r.Body)
// 	ubody := &defs.UserCredential{}
// 	if err := json.Unmarshal(res, ubody); err != nil {
// 		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
// 		return
// 	}
// 	pwd, err := dbops.GetUSerCredential(ubody.Username)

// 	if err != nil {
// 		sendErrorResponse(w, defs.ErrorNotAuthUser) //查询不到用户
// 		return
// 	}
// 	if ubody.Pwd != pwd {
// 		sendErrorResponse(w, defs.ErrorPwd) //密码错误
// 		return
// 	}
// 	resSessionId := session.GenerateNewSessionId(ubody.Username)
// 	si := &defs.SignedUp{Success: true, SessionId: resSessionId}
// 	if resp, err := json.Marshal(si); err != nil {
// 		sendErrorResponse(w, defs.ErrorInternalFaults)
// 		return
// 	} else {
// 		sendNormalResponse(w, string(resp), 200)
// 	}

// 	uname := p.ByName("user_name")
// 	io.WriteString(w, uname)
// }

func DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		//io.WriteString(w, "wrong")
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	uname := p.ByName("user_name") //验证用户是否对应
	if ubody.Username != uname {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	//验证用户密码
	userPwd, err := dbops.GetUSerCredential(uname)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	if userPwd != ubody.Pwd {
		sendErrorResponse(w, defs.ErrorPwd)
		return
	}
	//用户名及密码都验证通过
	err = dbops.DeleteUser(ubody.Username, ubody.Pwd)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	} else {
		sendNormalResponse(w, "ok", 202)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		//io.WriteString(w, "wrong")
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	// Validate the request body
	uname := p.ByName("user_name") //这一行是什么意思
	log.Printf("Login url name: %s", uname)

	log.Printf("Login body name: %s", ubody.Username)
	if uname != ubody.Username {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	log.Printf("%s", ubody.Username)
	pwd, err := dbops.GetUSerCredential(ubody.Username)
	log.Printf("Login pwd: %s", pwd)
	log.Printf("Login body pwd: %s", ubody.Pwd)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	si := &defs.SignedUp{Success: true, SessionId: id}
	if resp, err := json.Marshal(si); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

	//io.WriteString(w, "signed in")
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//首先校验用户
	if !ValidateUser(w, r) {
		log.Printf("GetUserInfo error: %v")
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	uname := p.ByName("user_name") //ByName获取的是url/后接的名字，需要将它与db中的name进行对比
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("GetUserCredential error: %v", err)
		return
	}
	ui := &defs.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	res, _ := ioutil.ReadAll(r.Body) //从request返回的应该是json字符串
	newVideo := &defs.NewVideo{}     //用于存储res中的json字符串的解析
	if err := json.Unmarshal(res, newVideo); err != nil {
		log.Printf("Unmarshal error: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}
	VI, err := dbops.AddNewVideo(newVideo.AuthorId, newVideo.Name) //VI存储的是video_info中的一条信息
	// type VideoInfo struct {
	// 	Id           string `json:"id"`
	// 	AuthorId     int    `json:"author_id"`
	// 	Name         string `jsonn:"name"`
	// 	DisplayCtime string `json:"display_ctime"`
	// }
	log.Printf("Author id : %d, name: %s \n", newVideo.AuthorId, newVideo.Name)
	if err != nil {
		log.Printf("dbops AddNewVideo error: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	if resp, err := json.Marshal(VI); err != nil {
		log.Printf("json.Marshal error: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
		return
	}

}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	uname := p.ByName("user_name")                                           //从request中得到User name
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec()) //GetCurrentTimestampSec() 获得当前时间
	if err != nil {
		log.Printf("ListVideoInfo error: %v", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	res := &defs.VideosInfo{Videos: vs}

	if resp, err := json.Marshal(res); err != nil {
		log.Printf("ListVideoInfo When Marshal error: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	vid := p.ByName("vid_id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	} else {
		go utils.SendDeleteVideoRequest(vid)
		sendNormalResponse(w, "", 204)
	}

	//204
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Printf("Error in PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
	} else {
		sendNormalResponse(w, "ok", 201)
	}

}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	cms := &defs.Comments{Comments: cm}
	if resp, err := json.Marshal(cms); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}
