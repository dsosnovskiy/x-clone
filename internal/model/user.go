package model

import (
	"time"
)

type User struct {
	UserID    int        `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username  string     `json:"username" gorm:"unique;not null"`
	Password  string     `json:"password" gorm:"not null"`
	FirstName string     `json:"first_name" gorm:"not null"`
	LastName  string     `json:"last_name" gorm:"not null"`
	Birthday  *time.Time `json:"birthday" gorm:"default:null"`
	Bio       *string    `json:"bio" gorm:"size:300;default:null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
