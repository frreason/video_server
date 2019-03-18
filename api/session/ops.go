package session

//sessionId 是用UUID算法生成的
import (
	"sync"
	"time"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
)

var sessionMap *sync.Map //用作缓存

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}
	r.Range(func(k, v interface{}) bool { //将从数据库拉出来的所有session存入map
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})

}

func GenerateNewSessionId(uname string) string { //最后要返回sessionId
	sid, _ := utils.NewUUID()
	ctime := nowInMilli()
	ttl := ctime + 30*60*1000                            //30min
	ss := &defs.SimpleSession{Username: uname, TTL: ttl} //这里取地址和不取地址有什么区别
	sessionMap.Store(sid, ss)                          //写入缓存
	dbops.InsertSession(sid, ttl, uname)                 //写入数据库
	return sid
}

func IsSessionExpired(sid string) (string, bool) { //判断session是否过期
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()
	if ok {
		if ss.(*defs.SimpleSession).TTL < ct { //ss.(*defs.SimpleSession).TTL是什么意思    过时
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*defs.SimpleSession).Username, false
	} else {
		ss, err := dbops.RetrieveSession(sid) //返回session对应的用户名及session
		if err != nil {
			return "", true
		}
		if ss.TTL < ct { //过时
			deleteExpiredSession(sid) //同时在缓存和数据库delete
			return "", true
		}
		sessionMap.Store(sid, ss)
		return ss.Username, false
	}
	return "", true
}
