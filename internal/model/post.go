package model

import (
	"time"
)

type Post struct {
	PostID         int       `json:"post_id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `json:"user_id" gorm:"index;not null;foreignKey:UserID"`
	Content        string    `json:"content" gorm:"size:1000;not null"`
	Likes          int       `json:"likes" gorm:"default:0"`
	Reposts        int       `json:"reposts" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	OriginalPostID *int      `json:"original_post_id" gorm:"index;default:null"`
	OriginalPost   *Post     `json:"original_post" gorm:"foreignKey:OriginalPostID;references:PostID"`
}

type Like struct {
	UserID      int `json:"user_id" gorm:"primaryKey;foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
	LikedPostID int `json:"liked_post_id" gorm:"primaryKey;foreignKey:LikedPostID;references:PostID;constraint:OnDelete:CASCADE"`
}
