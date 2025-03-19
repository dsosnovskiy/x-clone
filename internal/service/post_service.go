package service

import (
	"fmt"
	"x-clone/internal/model"
	"x-clone/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepository
	userRepo *repository.UserRepository
}

func NewPostService(postRepo *repository.PostRepository, userRepo *repository.UserRepository) *PostService {
	return &PostService{postRepo: postRepo, userRepo: userRepo}
}

func (s *PostService) CreatePost(post *model.Post, username string) error {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return err
	}

	post.UserID = user.UserID

	if err := s.postRepo.CreatePost(post); err != nil {
		return fmt.Errorf("failed to create post: %v", err)
	}
	return nil
}

func (s *PostService) GetUserPosts(username string) ([]model.Post, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return s.postRepo.GetUserPosts(user.UserID)
}

func (s *PostService) GetUserPostByID(userID, postID int) (*model.Post, error) {
	post, err := s.postRepo.GetUserPostByID(userID, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}
	return post, err
}

func (s *PostService) UpdatePostContentByID(postID, userID int, content string) error {
	if content == "" {
		return fmt.Errorf("post content cannot be empty")
	}

	postExists, err := s.postRepo.PostExists(postID)
	if err != nil {
		return err
	}
	if !postExists {
		return fmt.Errorf("post not found")
	}

	isOwner, err := s.postRepo.IsPostOwner(postID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("you are not the owner of this post")
	}

	if err := s.postRepo.UpdatePostContentByID(postID, content); err != nil {
		return err
	}
	return nil
}

func (s *PostService) DeletePostByID(postID, userID int) error {
	postExists, err := s.postRepo.PostExists(postID)
	if err != nil {
		return err
	}
	if !postExists {
		return fmt.Errorf("post not found")
	}

	isOwner, err := s.postRepo.IsPostOwner(postID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("you are not the owner of this post")
	}

	if err := s.postRepo.DeletePostByID(postID); err != nil {
		return err
	}
	return nil
}
