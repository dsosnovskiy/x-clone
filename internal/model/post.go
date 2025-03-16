package model

import (
	"time"
)

type Post struct {
	PostID         int       `json:"post_id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `json:"user_id" gorm:"index;not null;foreignKey:UserID"`
	Content        string    `json:"content" gorm:"size:1000;not null"`
	OriginalPostID *int      `json:"original_post_id" gorm:"index;default:null"`
	OriginalPost   *Post     `gorm:"foreignKey:OriginalPostID;references:PostID"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
