package controler

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"serve_video/api/dbops"
	"serve_video/api/model"
	"serve_video/api/session"
	"serve_video/api/utils"
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

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &model.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		//io.WriteString(w, "wrong")
		sendErrorRespinse(w, model.ErrorRequestBodyParseFailed)
		return
	}
	uname := p.ByName("username")
	log.Printf("Login url name: %s", uname)
	log.Printf("Login body name: %s", ubody.Username)
	if uname != ubody.Username {
		sendErrorRespinse(w, model.ErrorNotAuthUser)
		return
	}

	log.Printf("%s", ubody.Username)
	pwd, err := dbops.GetUserCreadential(ubody.Username)
	log.Printf("Login pwd: %s", pwd)
	log.Printf("Login body pwd: %s", ubody.Pwd)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		sendErrorRespinse(w, model.ErrorNotAuthUser)
		return
	}
	id := session.GenerateNewSessionId(ubody.Username)
	si := &model.SignedUp{Success: true, SessionId: id}
	if resp, err := json.Marshal(si); err != nil {
		sendErrorRespinse(w, model.ErrorInternalFailed)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

	//io.WriteString(w, uname)

}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		log.Printf("Unauthorized user\n")
		return
	}

	uname := p.ByName("username")
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("Error in GetUserInfo: %s", err)
		sendErrorRespinse(w, model.ErrorDBError)
		return
	}
	ui := &model.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorRespinse(w, model.ErrorInternalFailed)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		log.Printf("Unathorized user \n")
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &model.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil {
		log.Printf("%s", err)
		sendErrorRespinse(w, model.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddNewViedo(nvbody.AuthorId, nvbody.Name)
	log.Printf("Author id : %d, name: %s \n", nvbody.AuthorId, nvbody.Name)
	if err != nil {
		log.Printf("Error in AddNewVideo: %s", err)
		sendErrorRespinse(w, model.ErrorDBError)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		sendErrorRespinse(w, model.ErrorInternalFailed)
	} else {
		sendNormalResponse(w, string(resp), 201)
	}

}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	uname := p.ByName("username")
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ListAllvideos: %s", err)
		sendErrorRespinse(w, model.ErrorDBError)
		return
	}

	vsi := &model.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		sendErrorRespinse(w, model.ErrorInternalFailed)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		log.Printf("Error in DeletVideo: %s", err)
		sendErrorRespinse(w, model.ErrorDBError)
		return
	}

	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w, "", 204)
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	cbody := &model.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil {
		log.Printf("%s", err)
		sendErrorRespinse(w, model.ErrorRequestBodyParseFailed)
		return
	}

	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Printf("Error in PostComment: %s", err)
		sendErrorRespinse(w, model.ErrorDBError)
	} else {
		sendNormalResponse(w, "ok", 201)
	}

}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
		sendErrorRespinse(w, model.ErrorDBError)
		return
	}

	cms := &model.Comments{Comments: cm}
	if resp, err := json.Marshal(cms); err != nil {
		sendErrorRespinse(w, model.ErrorInternalFailed)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}
