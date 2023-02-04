package handlers

import (
	"encoding/base64"
	"errors"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"Forum"
)

var TemplateCreatePost = "templates/createPost.html"

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createPost" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}
	temp, err := template.ParseFiles(TemplateCreatePost)
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	user := r.Context().Value(ctxKeyUser).(Forum.User)
	switch r.Method {
	case http.MethodGet:
		if err := temp.Execute(w, user); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
	case http.MethodPost:
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		title, ok1 := r.Form["title"]
		text, ok2 := r.Form["text"]
		tag, ok3 := r.Form["tag"]
		if !ok1 || !ok2 || !ok3 {
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		if len(strings.TrimSpace(text[0])) == 0 || len(strings.TrimSpace(title[0])) == 0 {
			http.Redirect(w, r, "/createPost", http.StatusSeeOther)
		}
		tags := strings.Join(tag, " ")
		// upload image
		fileheader := r.MultipartForm.File["image"]
		for _, file := range fileheader {
			if file.Size > (21 << 20) {
				h.HandleErrorPage(w, http.StatusBadRequest, errors.New("image too large, you can upload files up to 20 MB"))
				return
			}
		}
		images, err := CreateImage(fileheader)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		// create post
		post := Forum.Post{
			UserId:     user.Id,
			Title:      title[0],
			Text:       text[0],
			Categories: tags,
			CreatedAt:  time.Now().Format("2 Jan 15:04:05"),
			Author:     user.Username,
			Image:      template.URL(images),
		}
		err = h.services.CreatePosts(post)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.HandleErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}
}

func CreateImage(fileheader []*multipart.FileHeader) (string, error) {
	var image string
	var images []string
	for _, f := range fileheader {
		file, err := f.Open()
		if err != nil {
			return "", err
		}
		filename := path.Base(f.Filename)
		if !strings.Contains(filename, ".jpg") && !strings.Contains(filename, "png") && !strings.Contains(filename, "gif") {
			return "", errors.New("you shou")
		}
		dest, err := os.Create("imageSource/" + filename)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(dest, file); err != nil {
			return "", err
		}
		// convert to  binary
		image, err = ConvertToBinary("imageSource/" + filename)
		if err != nil {
			return "", err
		}
		images = append(images, image)
	}
	img := strings.Join(images, " ")

	return img, nil
}

func ConvertToBinary(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	var base64Encoding string
	mimeType := http.DetectContentType(bytes)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	case "image/gif":
		base64Encoding += "data:image/gif;base64,"
	}
	base64Encoding += toBase64(bytes)
	return base64Encoding, nil
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
