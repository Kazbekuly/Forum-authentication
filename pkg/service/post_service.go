package service

import (
	"Forum"
	"Forum/pkg/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePosts(post Forum.Post) error {
	return s.repo.CreatePosts(post)
}

func (s *PostService) GetPostByUserID(id int) ([]Forum.Post, error) {
	return s.repo.GetPostByUserID(id)
}

func (s *PostService) GetPost() (**[]Forum.Post, error) {
	result, _ := s.repo.GetPost()
	return &result, nil
}

func (s *PostService) GetPostByTag(tags string) (**[]Forum.Post, error) {
	result, _ := s.repo.GetPostByTag(tags)
	return &result, nil
}

func (s *PostService) GetPostByPostID(postId int) (Forum.Post, error) {
	return s.repo.GetPostByPostID(postId)
}

func (s *PostService) UpdatePost(like string, dislike string, postId int) error {
	return s.repo.UpdatePost(like, dislike, postId)
}

func (s *PostService) LikedPosts(userId int) ([]Forum.Post, error) {
	return s.repo.LikedPosts(userId)
}
