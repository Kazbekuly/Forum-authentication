package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
)

type errorHTTP struct {
	Status  int
	Message error
}

var TemplateError = "templates/error.html"

func (h *Handler) HandleErrorPage(w http.ResponseWriter, status int, serverErr error) {
	if status >= http.StatusInternalServerError {
		log.Printf("something went wrong %s", serverErr)
	}
	errHTTP := errorHTTP{
		Status:  status,
		Message: serverErr,
	}
	w.WriteHeader(status)

	tmp, err := template.ParseFiles(TemplateError)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := tmp.Execute(w, errHTTP); err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
	}
}
