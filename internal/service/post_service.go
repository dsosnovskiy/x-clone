package service

import (
	"errors"
	"x-clone/internal/model"
	"x-clone/internal/repository"

	"gorm.io/gorm"
)

type PostService struct {
	postRepo *repository.PostRepository
	userRepo *repository.UserRepository
}

func NewPostService(postRepo *repository.PostRepository, userRepo *repository.UserRepository) *PostService {
	return &PostService{postRepo: postRepo, userRepo: userRepo}
}

func (s *PostService) CreatePost(post *model.Post) (*model.Post, error) {
	newPost, err := s.postRepo.CreatePost(post)
	if err != nil {
		return nil, errors.New("failed to create post")
	}
	return newPost, nil
}

func (s *PostService) GetUserPosts(userID int) (*[]model.Post, error) {
	posts, err := s.postRepo.GetUserPosts(userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetUserPostByID(userID, postID int) (*model.Post, error) {
	post, err := s.postRepo.GetUserPostByID(userID, postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("post not found")
		} else {
			return nil, err
		}
	}
	return post, nil
}

func (s *PostService) UpdatePostContentByID(userID, postID int, content string) (*model.Post, error) {
	updatedPost, err := s.postRepo.UpdatePostContentByID(userID, postID, content)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("post not found")
		} else {
			return nil, err
		}
	}
	return updatedPost, err
}

func (s *PostService) DeletePostByID(userID, postID int) error {
	if err := s.postRepo.DeletePostByID(userID, postID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("post not found")
		} else {
			return err
		}
	}
	return nil
}

func (s *PostService) LikePost(userID, postID int) error {
	if err := s.postRepo.LikePost(userID, postID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("post not found")
		} else {
			return err
		}
	}
	return nil
}

func (s *PostService) UnlikePost(userID, postID int) error {
	if err := s.postRepo.UnlikePost(userID, postID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("post not found")
		} else {
			return err
		}
	}
	return nil
}

func (s *PostService) RepostPost(userID, postID int) error {
	if err := s.postRepo.RepostPost(userID, postID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("post not found")
		} else {
			return err
		}
	}
	return nil
}

func (s *PostService) UndoRepostPost(userID, postID int) error {
	if err := s.postRepo.UndoRepostPost(userID, postID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("post not found")
		} else {
			return err
		}
	}
	return nil
}

func (s *PostService) GetUserReposts(userID int) (*[]model.Post, error) {
	reposts, err := s.postRepo.GetUserReposts(userID)
	if err != nil {
		return nil, err
	}
	return reposts, nil
}

func (s *PostService) QuotePost(userID, postID int, content string) (*model.Post, error) {
	post, err := s.postRepo.QuotePost(userID, postID, content)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("post not found")
		} else {
			return nil, err
		}
	}
	return post, nil
}
