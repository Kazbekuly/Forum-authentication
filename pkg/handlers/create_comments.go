package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Forum"
)

func (h *Handler) Comments(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}
	temp, err := template.ParseFiles(TemplateMyPost, TemplateCategory, TemplateHome)
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
		if err := r.ParseForm(); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		postId, ok1 := r.Form["postId"]
		comment, ok2 := r.Form["comment"]
		if !ok1 || !ok2 {
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		if len(strings.TrimSpace(comment[0])) == 0 {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			return
		}
		postID, err := strconv.Atoi(postId[0])
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		if user.Id == 0 {
			http.Redirect(w, r, "/signIn", 301)
			return
		}
		comments := Forum.Comments{
			UserId:    user.Id,
			PostId:    postID,
			Comment:   comment[0],
			Author:    user.Username,
			CreatedAt: time.Now().Format("2 Jan 15:04:05"),
		}
		err = h.services.AddComment(comments)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	default:
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
}
