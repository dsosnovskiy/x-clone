package repository

import (
	"fmt"
	"x-clone/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *model.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetUserPosts(userID int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetUserPostByID(userID, postID int) (*model.Post, error) {
	var post model.Post
	if err := r.db.Model(&model.Post{}).Where("user_id = ? AND post_id = ?", userID, postID).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) PostExists(postID int) (bool, error) {
	var post model.Post
	if err := r.db.Where("post_id = ?", postID).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *PostRepository) IsPostOwner(postID int, userID int) (bool, error) {
	var post model.Post
	if err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *PostRepository) UpdatePostContentByID(postID int, content string) error {
	var post model.Post
	if err := r.db.Model(&model.Post{}).Where("post_id = ?", postID).First(&post).Error; err != nil {
		return err
	}
	if err := r.db.Model(&post).Update("content", content).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeletePostByID(postID int) error {
	var post model.Post
	if err := r.db.Model(&model.Post{}).Where("post_id = ?", postID).First(&post).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&post).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) LikePost(userID, postID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// IsLiking
		var existingLike model.Like
		if err := tx.Where("user_id = ? AND liked_post_id = ?", userID, postID).First(&existingLike).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Continuing transaction...
			} else {
				return err
			}
		} else {
			return fmt.Errorf("you've already liked this post")
		}

		// CreateLike
		like := &model.Like{
			UserID:      userID,
			LikedPostID: postID,
		}
		if err := tx.Create(like).Error; err != nil {
			return err
		}

		// IncrementLikes
		if err := tx.Model(&model.Post{}).Where("post_id = ?", postID).Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostRepository) UnlikePost(userID, postID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// IsLiking
		var existingLike model.Like
		if err := tx.Where("user_id = ? AND liked_post_id = ?", userID, postID).First(&existingLike).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("you already don't like this post")
			}
			return err
		}

		// DeleteLike
		if err := tx.Delete(&existingLike).Error; err != nil {
			return err
		}

		// DecrementLikes
		if err := tx.Model(&model.Post{}).Where("post_id = ?", postID).Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}
