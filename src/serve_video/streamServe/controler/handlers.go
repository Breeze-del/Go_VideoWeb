package controler

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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
		log.Print("Error open video file")
		SendErrorResponse(w, 500, "Interal error"+err.Error())
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, req, "", time.Now(), video)
	defer video.Close()
}

func UploadHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	req.Body = http.MaxBytesReader(w, req.Body, modle.MAX_UPLOAD_SIZE)
	if err := req.ParseMultipartForm(modle.MAX_UPLOAD_SIZE); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}
	file, _, err := req.FormFile("file")
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Print("Read file error")
		SendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}
	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(modle.VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Print("write file error")
		SendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "upload successfully")
}

func TestPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")
	t.Execute(w, nil)
}
