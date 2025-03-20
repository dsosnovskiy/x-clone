package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	PostID         int       `json:"post_id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `json:"user_id" gorm:"index;not null;foreignKey:UserID"`
	Content        string    `json:"content" gorm:"size:1000;not null" validate:"required,min=10,max=1000"`
	OriginalPostID *int      `json:"original_post_id" gorm:"index;default:null" validate:"omitempty"`
	OriginalPost   *Post     `gorm:"foreignKey:OriginalPostID;references:PostID" validate:"omitempty"`
	Likes          int       `json:"likes" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Like struct {
	UserID      int `json:"user_id" gorm:"primaryKey;foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
	LikedPostID int `json:"liked_post_id" gorm:"primaryKey;foreignKey:LikedPostID;references:PostID;constraint:OnDelete:CASCADE"`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
