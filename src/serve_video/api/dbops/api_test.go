package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

// 删除表中所有记录,自增从1开始
func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("test", "123")
	if err != nil {
		t.Errorf("Error of AddUser : %v", err)
	}
}
func testGetUser(t *testing.T) {
	pwd, err := GetUserCreadential("test")
	if pwd != "123" && err != nil {
		t.Errorf("Error of GetUser : %v", err)
	}
}
func testDeleteUser(t *testing.T) {
	err := DeleteUserCreadential("test", "123")
	if err != nil {
		t.Errorf("Error of DelUser : %v", err)
	}
}
func testRegetUser(t *testing.T) {
	pwd, err := GetUserCreadential("test")
	if err != nil {
		t.Errorf("Error of RegetUser : %v", err)
	}
	if pwd != "" {
		t.Error("Error of RegetUser Failed")
	}
}

var tempvid string // videoId 全局
func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewViedo(1, "my-video")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	vi, err := GetVideoInfo(tempvid)
	if err != nil || vi != nil {
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

func TestCommentsFlow(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddComments", testAddComments)
	t.Run("ListComments", testListComments)
}
func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"
	err := AddNewComments(vid, aid, content)
	if err != nil {
		t.Errorf("addComments err : %s", err)
	}
}

func testListComments(t *testing.T) {
	vid := "12345"
	from := 1514764800
	// 获取当前本地时间戳，然后转成int类型
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("listComments err : %s", err)
	}
	for i, ele := range res {
		fmt.Printf("comment : %d, %v \n", i, ele)
	}
}
