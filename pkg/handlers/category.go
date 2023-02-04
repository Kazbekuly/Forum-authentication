package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"Forum"
)

var TemplateCategory = "templates/categories.html"

func (h *Handler) CategoryPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/categoryPost" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}
	temp, err := template.ParseFiles(TemplateCategory)
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	var result []Forum.Post
	switch r.Method {
	case http.MethodGet:
		if err := temp.Execute(w, result); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		tag, ok := r.Form["tag"]
		if !ok {
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		tags := strings.Join(tag, " ")
		category, err := h.services.GetPostByTag(tags)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		result = **category
		if err := temp.Execute(w, result); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
	default:
		h.HandleErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}
}
