package controler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"serve_video/streamServe/modle"
	"time"
)

func StreamHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := modle.VIDEO_DIR + vid
	video, err := os.Open(vl)
	if err != nil {
		SendErrorResponse(w, 500, "Interal error"+err.Error())
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, req, "", time.Now(), video)
	defer video.Close()
}

func UploadHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

}
