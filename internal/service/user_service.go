package service

import (
	"fmt"
	"x-clone/internal/model"
	"x-clone/internal/repository"
	"x-clone/pkg/hash"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *model.User) error {
	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = hashedPassword

	if err := s.userRepo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (s *UserService) FindUserByUsername(username string) (*model.User, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
