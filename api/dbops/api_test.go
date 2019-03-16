package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tempVid string

func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("mzj", "061365404abc")
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUSerCredential("mzj")
	if pwd != "061365404abc" || err != nil {
		t.Errorf("Error of GetUser: %v", err)
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("mzj", "061365404abc")
	if err != nil {
		t.Errorf("Error of deleteUser: %v", err)
	}
}

func testRegetUser(t *testing.T) {
	pwd, err := GetUSerCredential("mzj")
	if err != nil {
		t.Errorf("Error of RegetUser: %v", err)
	}
	if pwd != "" {
		t.Errorf("Delete User failed")
	}
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("Prepare add user", testAddUser)
	t.Run("Add new video", testAddVideoInfo)
	t.Run("Get video info", testGetVideoInfo)
	t.Run("Delete video info", testDeleteVideoInfo)
	t.Run("Reget video info", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	info, err := AddNewVideo(1, "first video")
	if err != nil {
		t.Errorf("Error of AddNewVideo: %v", err)
	}
	//fmt.Println(info)
	tempVid = info.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempVid)
	if err != nil {
		t.Errorf("Get video info error: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempVid)
	if err != nil {
		t.Errorf("Delete Video Info error: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	info, err := GetVideoInfo(tempVid)
	if err != nil || info != nil {
		//fmt.Println(sql.ErrNoRows)
		t.Errorf("Reget video info error: %v", err)
	}
}

func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddComments", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "hahahahaha"
	err := AddNewComments(vid, aid, content)
	err = AddNewComments(vid, aid, "aaaaa")
	err = AddNewComments(vid, aid, "bbbbb")
	err = AddNewComments(vid, aid, "ccccc")
	if err != nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

func testListComments(t *testing.T) { //不是很理解
	vid := "12345"
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))

	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}

	for i, ele := range res {
		fmt.Printf("comment: %d, %v \n", i, ele)
	}
}
