package handlers

import (
	"errors"
	"html/template"
	"net/http"

	"Forum"
)

var TemplateMyPost = "templates/myposts.html"

func (h *Handler) MyPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myposts" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}
	temp, err := template.ParseFiles(TemplateMyPost)
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	user := r.Context().Value(ctxKeyUser).(Forum.User)
	switch r.Method {
	case http.MethodGet:
		post, err := h.services.GetPostByUserID(user.Id)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		result := Forum.MyPost{
			User:  user.Username,
			Posts: post,
		}
		if err := temp.Execute(w, result); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
	case http.MethodPost:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.HandleErrorPage(w, http.StatusMethodNotAllowed, nil)
	}
}
