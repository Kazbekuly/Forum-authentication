package repository

import (
	"database/sql"

	"Forum"
)

type Authorization interface {
	CreateUser(user Forum.User) error
	GetUser(username, email string) (Forum.User, error)
	AddToken(token Forum.Token) (Forum.Token, error)
	CheckInvalid(email, username string) (Forum.User, error)
	GetToken(token string) (Forum.Token, error)
	GetUserByToken(token string) (Forum.User, error)
	DeleteToken(token string) error
	DeleteTokenByUserID(id int) error
}

type Post interface {
	CreatePosts(post Forum.Post) error
	GetPostByUserID(id int) ([]Forum.Post, error)
	GetPost() (*[]Forum.Post, error)
	GetPostByTag(tags string) (*[]Forum.Post, error)
	GetPostByPostID(postId int) (Forum.Post, error)
	UpdatePost(like string, dislike string, postId int) error
	LikedPosts(userId int) ([]Forum.Post, error)
}

type Comment interface {
	AddComments(comments Forum.Comments) error
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

type Repository struct {
	Authorization
	Post
	Comment
	Evaluate
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSQL(db),
		Post:          NewPostSQL(db),
		Comment:       NewCommentSQL(db),
		Evaluate:      NewEvaluateSQL(db),
	}
}
