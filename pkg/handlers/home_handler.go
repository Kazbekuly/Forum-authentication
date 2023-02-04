package handlers

import (
	"errors"
	"html/template"
	"net/http"

	"Forum"
)

var TemplateHome = "templates/home.html"

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}
	if r.Method != http.MethodGet {
		h.HandleErrorPage(w, http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	user := r.Context().Value(ctxKeyUser).(Forum.User)
	post, err := h.services.GetPost()
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	result := Forum.MyPost{
		User:  user.Username,
		Posts: **post,
	}
	tmpl, err := template.ParseFiles(TemplateHome)
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	if err := tmpl.Execute(w, result); err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
}
