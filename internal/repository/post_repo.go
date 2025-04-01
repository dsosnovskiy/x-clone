package repository

import (
	"errors"
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

func (r *PostRepository) GetUserPosts(userID int) (*[]model.Post, error) {
	var posts []model.Post
	if err := r.db.Preload("OriginalPost", func(db *gorm.DB) *gorm.DB {
		return db.Preload("OriginalPost") // Recursive preload OriginalPost *Post
	}).Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return &posts, nil
}

func (r *PostRepository) GetUserPostByID(userID, postID int) (*model.Post, error) {
	var post model.Post
	if err := r.db.Preload("OriginalPost", func(db *gorm.DB) *gorm.DB {
		return db.Preload("OriginalPost") // Recursive preload OriginalPost *Post
	}).Where("user_id = ? AND post_id = ?", userID, postID).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) UpdatePostContentByID(userID, postID int, content string) error {
	if err := r.db.Model(&model.Post{}).Where("user_id = ? AND post_id = ?", userID, postID).Update("content", content).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeletePostByID(userID, postID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Deleting all reposts associated with this post
		if err := tx.Where("reposted_post_id = ?", postID).Delete(&model.Repost{}).Error; err != nil {
			return err
		}

		// Deleting all likes associated with this post
		if err := tx.Where("liked_post_id = ?", postID).Delete(&model.Like{}).Error; err != nil {
			return err
		}

		// Deleting all quotes asscodiated with this post
		var quotes []model.Post
		if err := tx.Where("original_post_id = ?", postID).Find(&quotes).Error; err != nil {
			return err
		}
		for _, quote := range quotes {
			if err := r.DeletePostByID(quote.UserID, quote.PostID); err != nil {
				return err
			}
		}

		// Deleting the post itself
		if err := tx.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Post{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostRepository) LikePost(userID, postID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// IsLiked
		var existingLike model.Like
		if err := tx.Where("user_id = ? AND liked_post_id = ?", userID, postID).First(&existingLike).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		} else {
			return errors.New("you've already liked this post")
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
		// IsLiked
		var existingLike model.Like
		if err := tx.Where("user_id = ? AND liked_post_id = ?", userID, postID).First(&existingLike).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("you already don't like this post")
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

func (r *PostRepository) RepostPost(userID, postID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// IsReposted
		var existingRepost model.Repost
		if err := tx.Where("user_id = ? AND reposted_post_id = ?", userID, postID).First(&existingRepost).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		} else {
			return errors.New("you've already reposted this post")
		}

		// CreateRepost
		repost := &model.Repost{
			UserID:         userID,
			RepostedPostID: postID,
		}
		if err := tx.Create(repost).Error; err != nil {
			return err
		}

		// IncrementReposts
		if err := tx.Model(&model.Post{}).Where("post_id = ?", postID).Update("reposts", gorm.Expr("reposts + 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostRepository) UndoRepostPost(userID, postID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// IsReposted
		var existingRepost model.Repost
		if err := tx.Where("user_id = ? AND reposted_post_id = ?", userID, postID).First(&existingRepost).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("you haven't reposted this post")
			}
			return err
		}

		// DeleteRepost
		if err := tx.Delete(&existingRepost).Error; err != nil {
			return err
		}

		// DecrementReposts
		if err := tx.Model(&model.Post{}).Where("post_id = ?", postID).Update("reposts", gorm.Expr("reposts - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostRepository) GetUserReposts(userID int) (*[]model.Post, error) {
	var reposts []model.Repost
	if err := r.db.Where("user_id = ?", userID).Find(&reposts).Error; err != nil {
		return nil, err
	}

	var postIDs []int
	for _, repost := range reposts {
		postIDs = append(postIDs, repost.RepostedPostID)
	}

	var posts []model.Post
	if err := r.db.Where("post_id IN ?", postIDs).Find(&posts).Error; err != nil {
		return nil, err
	}

	return &posts, nil
}

func (r *PostRepository) QuotePost(userID, postID int, content string) (*model.Post, error) {
	post := &model.Post{
		UserID:         userID,
		Content:        content,
		OriginalPostID: &postID,
	}

	if err := r.db.Create(post).Error; err != nil {
		return nil, err
	}

	var quotedPost model.Post
	if err := r.db.Preload("OriginalPost", func(db *gorm.DB) *gorm.DB {
		return db.Preload("OriginalPost") // Recursive preload OriginalPost *Post
	}).Where("post_id = ?", post.PostID).First(&quotedPost).Error; err != nil {
		return nil, err
	}

	return &quotedPost, nil
}
