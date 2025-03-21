package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	PostID         int       `json:"post_id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `json:"user_id" gorm:"index;not null;foreignKey:UserID"`
	Content        string    `json:"content" gorm:"size:1000;not null" validate:"required,min=1,max=1000"`
	Likes          int       `json:"likes" gorm:"default:0"`
	Reposts        int       `json:"reposts" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	OriginalPostID *int      `json:"original_post_id" gorm:"index;default:null"`
	OriginalPost   *Post     `gorm:"foreignKey:OriginalPostID;references:PostID"`
}

type Like struct {
	UserID      int `json:"user_id" gorm:"primaryKey;foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
	LikedPostID int `json:"liked_post_id" gorm:"primaryKey;foreignKey:LikedPostID;references:PostID;constraint:OnDelete:CASCADE"`
}

func (p *Post) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

func ValidateContent(content string) error {
	validate := validator.New()
	type Content struct {
		Content string `validate:"required,min=1,max=1000"`
	}
	contentData := Content{
		Content: content,
	}
	if err := validate.Struct(contentData); err != nil {
		return err
	}
	return nil
}
