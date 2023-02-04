package handlers

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"Forum"
)

func (h *Handler) LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likeComment" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
		return
	}
	temp, err := template.ParseFiles(TemplateMyPost, TemplateCategory, TemplateHome, TemplateLikedPosts)
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
		commentId, ok1 := r.Form["commentId"]
		if !ok1 {
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		id, err := strconv.Atoi(commentId[0])
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		if user.Id == 0 {
			http.Redirect(w, r, "/signIn", 301)
			return
		}
		evaluates := Forum.EvaluateComment{
			CommentId: id,
			UserId:    user.Id,
			Vote:      1,
		}
		reaction, err := h.services.CheckUserComment(user.Id, id)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		if reaction.UserId != user.Id && reaction.CommentId != id {
			if err := h.services.CreateEvaluateComment(evaluates); err != nil {
				h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
				return
			}
		}
		if reaction.UserId == user.Id && reaction.CommentId == id && reaction.Vote != 1 {
			err = h.services.UpdateCommentVote(user.Id, id, 1)
			if err != nil {
				h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
				return
			}
		}
		if reaction.UserId == user.Id && reaction.CommentId == id && reaction.Vote == 1 {
			if err = h.services.CheckCommentVote(user.Id, id, 1); err != nil {
				h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
				return
			}
		}
		count, err := h.services.EvaluateCommentCount(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
				return
			}
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		err = h.services.UpdateComment(count.Like, count.Dislike, id)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	}
}

func (h *Handler) DislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dislikeComment" {
		h.HandleErrorPage(w, http.StatusNotFound, errors.New("error with url"))
		return
	}
	temp, err := template.ParseFiles(TemplateMyPost, TemplateCategory, TemplateHome, TemplateLikedPosts)
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
		commentId, ok1 := r.Form["commentId"]
		if !ok1 {
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
		id, err := strconv.Atoi(commentId[0])
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		if user.Id == 0 {
			http.Redirect(w, r, "/signIn", 301)
			return
		}
		evaluates := Forum.EvaluateComment{
			CommentId: id,
			UserId:    user.Id,
			Vote:      -1,
		}
		reaction, err := h.services.CheckUserComment(user.Id, id)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		if reaction.UserId != user.Id && reaction.CommentId != id {
			if err := h.services.CreateEvaluateComment(evaluates); err != nil {
				h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
				return
			}
		}
		if reaction.UserId == user.Id && reaction.CommentId == id && reaction.Vote != -1 {
			err = h.services.UpdateCommentVote(user.Id, id, -1)
			if err != nil {
				h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
				return
			}
		}
		if reaction.UserId == user.Id && reaction.CommentId == id && reaction.Vote == -1 {
			if err = h.services.CheckCommentVote(user.Id, id, -1); err != nil {
				h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
				return
			}
		}
		count, err := h.services.EvaluateCommentCount(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.HandleErrorPage(w, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
				return
			}
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		err = h.services.UpdateComment(count.Like, count.Dislike, id)
		if err != nil {
			h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
			return
		}
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	}
}
