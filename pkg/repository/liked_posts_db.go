package repository

import (
	"Forum"
	"context"
	"time"
)

func (r *PostSQL) LikedPosts(userId int) ([]Forum.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var evaluate Forum.Evaluate
	var post Forum.Post
	var posts []int
	var result []Forum.Post
	row, err := r.db.QueryContext(ctx, "SELECT PostId FROM evaluate WHERE UserId=$1 AND Vote=$2 ", userId, "1")
	if err != nil {
		return nil, err
	}
	for row.Next() {
		if err := row.Scan(&evaluate.PostId); err != nil {
			return nil, err
		}
		posts = append(posts, evaluate.PostId)
	}
	for _, id := range posts {
		if err := r.db.QueryRowContext(ctx, "SELECT * FROM post WHERE id=$1", id).
			Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Categories, &post.CreatedAt, &post.Author, &post.Like, &post.Dislike, &post.Image); err != nil {
			return nil, err
		}
		result = append(result, post)
	}
	return result, nil
}
