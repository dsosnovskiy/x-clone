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
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
