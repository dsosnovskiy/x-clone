package repository

import (
	"errors"
	"x-clone/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(userID int) (*model.User, error) {
	var user model.User
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FollowUser(followerID, followingID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// IsFollowing
		var existingFollower model.Follower
		if err := tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existingFollower).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Continuing transaction...
			} else {
				return err
			}
		} else {
			return errors.New("you are already following this user")
		}

		// CreateFollower
		follower := &model.Follower{
			FollowerID:  followerID,
			FollowingID: followingID,
		}
		if err := tx.Create(follower).Error; err != nil {
			return err
		}

		// IncrementFollowers
		if err := tx.Model(&model.User{}).Where("user_id = ?", followingID).Update("followers", gorm.Expr("followers + 1")).Error; err != nil {
			return err
		}
		// IncrementFollowing
		if err := tx.Model(&model.User{}).Where("user_id = ?", followerID).Update("following", gorm.Expr("following + 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) StopFollowingUser(followerID, followingID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// IsFollowing
		var existingFollower model.Follower
		if err := tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existingFollower).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("you are not following this user")
			}
			return err
		}

		// DeleteFollower
		if err := tx.Delete(&existingFollower).Error; err != nil {
			return err
		}

		// DecrementFollowers
		if err := tx.Model(&model.User{}).Where("user_id = ?", followingID).Update("followers", gorm.Expr("followers - 1")).Error; err != nil {
			return err
		}
		// DecrementFollowing
		if err := tx.Model(&model.User{}).Where("user_id = ?", followerID).Update("following", gorm.Expr("following - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) GetFollowersByUser(userID int) ([]model.User, error) {
	var user model.User

	err := r.db.
		Preload("FollowersList").
		Where("user_id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return user.FollowersList, nil
}

func (r *UserRepository) GetFollowingByUser(userID int) ([]model.User, error) {
	var user model.User

	err := r.db.
		Preload("FollowingList").
		Where("user_id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return user.FollowingList, nil
}

func (r *UserRepository) ChangeProfile(userID int, username, firstName, lastName, birthday, bio *string) error {
	if username != nil {
		if err := r.db.Model(&model.User{}).Where("user_id = ?", userID).Update("username", username).Error; err != nil {
			return err
		}
	}
	if firstName != nil {
		if err := r.db.Model(&model.User{}).Where("user_id = ?", userID).Update("first_name", firstName).Error; err != nil {
			return err
		}
	}
	if lastName != nil {
		if err := r.db.Model(&model.User{}).Where("user_id = ?", userID).Update("last_name", lastName).Error; err != nil {
			return err
		}
	}
	if birthday != nil {
		if err := r.db.Model(&model.User{}).Where("user_id = ?", userID).Update("birthday", birthday).Error; err != nil {
			return err
		}
	}
	if bio != nil {
		if err := r.db.Model(&model.User{}).Where("user_id = ?", userID).Update("bio", bio).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) ChangePassword(userID int, hashedNewPassword string) error {
	if err := r.db.Model(&model.User{}).Where("user_id = ?", userID).Update("password", hashedNewPassword).Error; err != nil {
		return err
	}
	return nil
}
