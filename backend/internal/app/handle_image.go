package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ImageResponse struct {
	ImageId string `json:"imageId"`
}

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

func (a *Application) handlePostImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
		a.badRequestResponse(w, r, err)
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", fileName))
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	defer dst.Close()

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	writeJsonResponse(w, http.StatusCreated, ImageResponse{ImageId: fileName}, nil)
}

func (a *Application) handleGetImage(w http.ResponseWriter, r *http.Request) {
	imageId := filepath.Clean(getField(r, 0))
	http.ServeFile(w, r, fmt.Sprintf("./uploads/%s", imageId))
}
