package model

import (
	"time"
)

type User struct {
	UserID        int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username      string    `json:"username" gorm:"unique;not null"`
	Password      string    `json:"password" gorm:"not null"`
	FirstName     string    `json:"first_name" gorm:"not null"`
	LastName      string    `json:"last_name" gorm:"not null"`
	Birthday      *string   `json:"birthday" gorm:"default:null"`
	Bio           *string   `json:"bio" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	Followers     int       `json:"followers" gorm:"default:0"`
	Following     int       `json:"following" gorm:"default:0"`
	FollowersList []User    `gorm:"many2many:followers;foreignKey:UserID;joinForeignKey:FollowingID;References:UserID;joinReferences:FollowerID"`
	FollowingList []User    `gorm:"many2many:followers;foreignKey:UserID;joinForeignKey:FollowerID;References:UserID;joinReferences:FollowingID"`
}

type Follower struct {
	FollowerID  int `json:"follower_id" gorm:"primaryKey;foreignKey:FollowerID;references:UserID;constraint:OnDelete:CASCADE"`
	FollowingID int `json:"following_id" gorm:"primaryKey;foreignKey:FollowingID;references:UserID;constraint:OnDelete:CASCADE"`
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
