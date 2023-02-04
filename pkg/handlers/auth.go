package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"time"

	"Forum"
	"Forum/pkg/service"
)

var (
	TemplateSignUp = "templates/signUp.html"
	TemplateSignIn = "templates/signIn.html"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signUp" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}
	temp, err := template.ParseFiles(TemplateSignUp)
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := temp.Execute(w, nil); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		email, ok1 := r.Form["email"]
		username, ok2 := r.Form["username"]
		password, ok3 := r.Form["password"]
		if !ok1 || !ok2 || !ok3 {
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		user := Forum.User{
			Email:    email[0],
			Username: username[0],
			Password: password[0],
		}
		err := h.services.Authorization.CreateUser(user)
		if err != nil {
			if errors.Is(err, service.ErrCheckInvalid) {
				h.HandleErrorPage(w, http.StatusBadRequest, service.ErrCheckInvalid)
				return
			} else if errors.Is(err, sql.ErrNoRows) {
				h.HandleErrorPage(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
				return
			}
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		http.Redirect(w, r, "/signIn", http.StatusSeeOther)
	default:
		h.HandleErrorPage(w, http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles(TemplateSignIn)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		if err := temp.Execute(w, nil); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		username, ok1 := r.Form["username"]
		password, ok2 := r.Form["password"]
		if !ok1 || !ok2 {
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		token, err := h.services.Authorization.GenerateToken(username[0], password[0], false)
		if err != nil {
			if errors.Is(err, service.ErrorWrongPassword) {
				h.HandleErrorPage(w, http.StatusBadRequest, service.ErrorWrongPassword)
				return
			}
			h.HandleErrorPage(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: token.AuthToken,
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.HandleErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}
	if r.Method != http.MethodGet {
		h.HandleErrorPage(w, http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.HandleErrorPage(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}
		h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}
	if err := h.services.DeleteToken(cookie.Value); err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Time{},
		Path:    "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
