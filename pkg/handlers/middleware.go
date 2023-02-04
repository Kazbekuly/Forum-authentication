package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"Forum"
)

const ctxKeyUser ctxKey = iota

type ctxKey int8

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var user Forum.User
		cookie, err := request.Cookie("session_token")
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				h.HandleErrorPage(writer, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
				return
			}
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, Forum.User{})))
			return
		}

		token, err := h.services.GetToken(cookie.Value)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				h.HandleErrorPage(writer, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
				return
			}
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, Forum.User{})))
			return
		}
		if token.ExpiresAT.Before(time.Now()) {
			h.services.DeleteToken(cookie.Value)
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, Forum.User{})))
			return
		}
		user, err = h.services.GetUserByToken(token.AuthToken)
		if err != nil {
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, Forum.User{})))
			return
		}
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyUser, user)))
	}
}
