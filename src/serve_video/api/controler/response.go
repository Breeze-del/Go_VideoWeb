package controler

import (
	"encoding/json"
	"io"
	"net/http"
	"serve_video/api/model"
)

// 处理所有HTTP返回消息
func sendErrorRespinse(w http.ResponseWriter, errResp model.ErroResponse) {
	// 返回错误信息
	w.WriteHeader(errResp.HttpSc)
	resStr, _ := json.Marshal(&errResp.Error)
	w.Write(resStr)
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
