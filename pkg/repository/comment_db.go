package repository

import (
	"Forum"
	"context"
	"database/sql"
	"time"
)

type CommentSQL struct {
	db *sql.DB
}

func NewCommentSQL(db *sql.DB) *CommentSQL {
	return &CommentSQL{db: db}
}

func (r *CommentSQL) AddComments(comments Forum.Comments) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	query, err := r.db.PrepareContext(ctx, "INSERT INTO comment (UserId, PostId, comment, createdAt, author, like_vote, dislike) values ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	_, err = query.Exec(comments.UserId, comments.PostId, comments.Comment, comments.CreatedAt, comments.Author, "0", "0")
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentSQL) GetCommentByPost(postId int) ([]Forum.Comments, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var result []Forum.Comments
	var comment Forum.Comments
	row, err := r.db.QueryContext(ctx, "SELECT * FROM comment WHERE PostId=$1", postId)
	if err != nil {
		return result, err
	}
	for row.Next() {
		if err := row.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.Comment, &comment.CreatedAt, &comment.Author, &comment.Like, &comment.Dislike); err != nil {
			return result, err
		}
		evaluates, err := NewEvaluateSQL(r.db).EvaluateCommentCount(comment.Id)
		if err != nil {
			return nil, err
		}
		comment.Like = evaluates.Like
		comment.Dislike = evaluates.Dislike
		result = append(result, comment)
	}
	return result, nil
}

func (r *CommentSQL) UpdateComment(like string, dislike string, commentId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	_, err := r.db.ExecContext(ctx, "UPDATE comment SET like_vote=$1, dislike=$2 WHERE id=$3 ", like, dislike, commentId)
	if err != nil {
		return err
	}
	return nil
}
