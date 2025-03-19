package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	UserID        int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username      string    `json:"username" gorm:"unique;not null" validate:"required,min=6,max=20"`
	Password      string    `json:"password" gorm:"not null" validate:"required,min=7,max=32"`
	FirstName     string    `json:"first_name" gorm:"not null" validate:"required,min=2,max=32"`
	LastName      string    `json:"last_name" gorm:"not null" validate:"required,min=2,max=32"`
	Birthday      *string   `json:"birthday" gorm:"default:null" validate:"omitempty,datetime=2006-01-02"`
	Bio           *string   `json:"bio" gorm:"size:300;default:null" validate:"max=300"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Followers     int       `json:"followers" gorm:"default:0"`
	Following     int       `json:"following" gorm:"default:0"`
	FollowersList []User    `gorm:"many2many:followers;foreignKey:UserID;joinForeignKey:FollowingID;References:UserID;joinReferences:FollowerID"`
	FollowingList []User    `gorm:"many2many:followers;foreignKey:UserID;joinForeignKey:FollowerID;References:UserID;joinReferences:FollowingID"`
}

type Follower struct {
	FollowerID  int `json:"follower_id" gorm:"primaryKey;foreignKey:FollowerID;references:UserID;constraint:OnDelete:CASCADE"`
	FollowingID int `json:"following_id" gorm:"primaryKey;foreignKey:FollowingID;references:UserID;constraint:OnDelete:CASCADE"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserResponse struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  *string   `json:"birthday"`
	Bio       *string   `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	Followers int       `json:"followers"`
	Following int       `json:"following"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		UserID:    u.UserID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Birthday:  u.Birthday,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
		Followers: u.Followers,
		Following: u.Following,
	}
}
