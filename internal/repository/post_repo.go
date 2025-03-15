package repository

import (
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

func (r *PostRepository) GetUserPosts(UserID int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.Where("user_id = ?", UserID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
