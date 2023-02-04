package repository

import (
	"context"
	"database/sql"
	"time"

	"Forum"
)

type PostSQL struct {
	db *sql.DB
}

func NewPostSQL(db *sql.DB) *PostSQL {
	return &PostSQL{db: db}
}

func (r *PostSQL) CreatePosts(post Forum.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	query, err := r.db.PrepareContext(ctx, "INSERT INTO post (UserId, title, text, category, createdAt, author, like_vote, dislike, image) values ($1, $2, $3, $4, $5, $6,$7,$8, $9)")
	if err != nil {
		return err
	}
	if _, err = query.Exec(post.UserId, post.Title, post.Text, post.Categories, post.CreatedAt, post.Author, "0", "0", post.Image); err != nil {
		return err
	}
	return nil
}

func (r *PostSQL) GetPostByUserID(id int) ([]Forum.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var result []Forum.Post
	var post Forum.Post
	row, err := r.db.QueryContext(ctx, "SELECT * FROM post WHERE UserID=$1 ", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Categories, &post.CreatedAt, &post.Author, &post.Like, &post.Dislike, &post.Image); err != nil {
			return result, err
		}
		evaluates, err := NewEvaluateSQL(r.db).EvaluateCount(post.Id)
		if err != nil {
			return nil, err
		}
		post.Like = evaluates.Like
		post.Dislike = evaluates.Dislike
		result = append(result, post)
	}
	return result, nil
}

func (r *PostSQL) GetPost() (*[]Forum.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var result []Forum.Post
	var post Forum.Post
	row, err := r.db.QueryContext(ctx, "SELECT *FROM post")
	if err != nil {
		return &result, err
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Categories, &post.CreatedAt, &post.Author, &post.Like, &post.Dislike, &post.Image); err != nil {
			return &result, err
		}
		evaluates, err := NewEvaluateSQL(r.db).EvaluateCount(post.Id)
		if err != nil {
			return &[]Forum.Post{}, err
		}
		post.Like = evaluates.Like
		post.Dislike = evaluates.Dislike
		result = append(result, post)
	}

	return &result, nil
}

func (r *PostSQL) GetPostByTag(tags string) (*[]Forum.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var result []Forum.Post
	var post Forum.Post
	row, err := r.db.QueryContext(ctx, "SELECT * FROM post WHERE category like "+"'"+"%"+tags+"%'")
	if err != nil {
		return &result, err
	}
	for row.Next() {
		if err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Categories, &post.CreatedAt, &post.Author, &post.Like, &post.Dislike, &post.Image); err != nil {
			return &result, err
		}
		evaluates, err := NewEvaluateSQL(r.db).EvaluateCount(post.Id)
		if err != nil {
			return &result, err
		}
		post.Like = evaluates.Like
		post.Dislike = evaluates.Dislike
		result = append(result, post)
	}
	return &result, nil
}

func (r *PostSQL) GetPostByPostID(postId int) (Forum.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var post Forum.Post
	row := r.db.QueryRowContext(ctx, "SELECT * FROM post WHERE id=$1", postId)
	if err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Categories, &post.CreatedAt, &post.Author, &post.Like, &post.Dislike, &post.Image); err != nil {
		return post, err
	}
	return post, nil
}

func (r *PostSQL) UpdatePost(like string, dislike string, postId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	_, err := r.db.ExecContext(ctx, "UPDATE post SET like_vote=$1, dislike=$2 WHERE id=$3 ", like, dislike, postId)
	if err != nil {
		return err
	}
	return nil
}
