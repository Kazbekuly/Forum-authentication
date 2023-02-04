package service

import (
	"Forum"
	"Forum/pkg/repository"
)

type Authorization interface {
	CreateUser(user Forum.User) error
	GenerateToken(username, password string, oauth bool) (Forum.Token, error)
	GetToken(token string) (Forum.Token, error)
	GetUserByToken(token string) (Forum.User, error)
	DeleteToken(token string) error
}

type Post interface {
	CreatePosts(post Forum.Post) error
	GetPostByUserID(id int) ([]Forum.Post, error)
	GetPost() (**[]Forum.Post, error)
	GetPostByTag(tags string) (**[]Forum.Post, error)
	GetPostByPostID(postId int) (Forum.Post, error)
	UpdatePost(like string, dislike string, postId int) error
	LikedPosts(userId int) ([]Forum.Post, error)
}

type Comment interface {
	AddComment(comment Forum.Comments) error
	GetCommentByPost(postId int) ([]Forum.Comments, error)
	UpdateComment(like string, dislike string, commentId int) error
}

type Evaluate interface {
	CreateEvaluates(postLike Forum.Evaluate) error
	CheckUserPost(userId, postId int) (Forum.Evaluate, error)
	UpdateVote(userId, postId, newVote int) error
	EvaluateCount(postId int) (Forum.Vote, error)
	CheckUserComment(userId, commentId int) (Forum.EvaluateComment, error)
	CreateEvaluateComment(commentVote Forum.EvaluateComment) error
	UpdateCommentVote(userId, commentId, newVote int) error
	EvaluateCommentCount(commentId int) (Forum.Vote, error)
	CheckVote(userId, postId, vote int) error
	CheckCommentVote(userId, postId, vote int) error
}

type Service struct {
	Authorization
	Post
	Comment
	Evaluate
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Post:          NewPostService(repo.Post),
		Comment:       NewCommentService(repo.Comment),
		Evaluate:      NewEvaluateService(repo.Evaluate),
	}
}
