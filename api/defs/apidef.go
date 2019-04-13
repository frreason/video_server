package defs

type User struct {
	Id        int //Id是什么意思  数据库里每个元组的Id
	LoginName string
	Pwd       string
}

type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type SignedUp struct {
	Success   bool   `json:"succsee"`
	SessionId string `json:"session_id"`
}

type VideoInfo struct {
	Id           string `json:"id"`
	AuthorId     int    `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type Comment struct {
	Id      string `json:"id"`
	VideoId string `json:"video_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

type SimpleSession struct {
	Username string
	TTL      int64
}

type UserInfo struct {
	Id int `json:"id"`
}

type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}

type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

type UserSession struct {
	Username  string `json:"user_name"`
	SessionId string `json:"session_id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type Comments struct {
	Comments []*Comment `json:"comments"`
}
