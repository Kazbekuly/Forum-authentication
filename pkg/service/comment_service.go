package service

import (
	"Forum"
	"Forum/pkg/repository"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) AddComment(comment Forum.Comments) error {
	return s.repo.AddComments(comment)
}

func (s *CommentService) GetCommentByPost(postId int) ([]Forum.Comments, error) {
	return s.repo.GetCommentByPost(postId)
}

func (s *CommentService) UpdateComment(like string, dislike string, commentId int) error {
	return s.repo.UpdateComment(like, dislike, commentId)
}
