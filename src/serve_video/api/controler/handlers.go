package controler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"serve_video/api/dbops"
	"serve_video/api/model"
	"serve_video/api/session"
)

// 注册用户
func CreateUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(req.Body)
	ubody := &model.UserCredential{}
	// 解析用户信息
	if err := json.Unmarshal(res, ubody); err != nil {
		// 如果解析失败，那么返回解析失败
		sendErrorRespinse(w, model.ErrorRequestBodyParseFailed)
		return
	}
	// 将用户信息写入数据库
	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorRespinse(w, model.ErrorDBError)
		return
	}
	// 产生session写入数据库和cache 返回sessionId
	id := session.GenerateNewSessionId(ubody.Username)
	// 将sessionId返回给前端
	su := &model.SignedUp{Success: true, SessionId: id}
	if resp, err := json.Marshal(su); err != nil {
		sendErrorRespinse(w, model.ErrorInternalFailed)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}

func Login(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	usname := p.ByName("user_name")
	io.WriteString(w, usname)
}
