package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(s interface{}) error {
	return validate.Struct(s)
}

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required,min=6,max=20"`
	Password  string  `json:"password" validate:"required,min=7,max=32"`
	FirstName string  `json:"first_name" validate:"required,min=2,max=32"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=32"`
	Birthday  *string `json:"birthday" validate:"omitempty,datetime=2006-01-02"`
	Bio       *string `json:"bio" validate:"omitempty,min=1,max=300"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=6,max=20"`
	Password string `json:"password" validate:"required,min=7,max=32"`
}

type ProfileUpdateRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=6,max=20"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=32"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=32"`
	Birthday  *string `json:"birthday" validate:"omitempty,datetime=2006-01-02"`
	Bio       *string `json:"bio" validate:"omitempty,min=1,max=300"`
}

type PasswordChangeRequest struct {
	OldPassword     string `json:"old_password" validate:"required,min=7,max=32"`
	NewPassword     string `json:"new_password" validate:"required,min=7,max=32,nefield=OldPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=7,max=32,eqfield=NewPassword"`
}

type ContentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}
