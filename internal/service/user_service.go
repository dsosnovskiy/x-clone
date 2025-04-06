package service

import (
	"errors"
	"x-clone/internal/model"
	"x-clone/internal/repository"
	"x-clone/pkg/utils/hash"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByUsername(username string) (*model.UserResponse, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	userResponse := user.ToResponse()
	return &userResponse, nil
}

func (s *UserService) FollowUser(followerID, followingID int) error {
	if followerID == followingID {
		return errors.New("you cannot follow yourself")
	}
	return s.userRepo.FollowUser(followerID, followingID)
}

func (s *UserService) StopFollowingUser(followerID, followingID int) error {
	if followerID == followingID {
		return errors.New("you cannot stop following yourself")
	}
	return s.userRepo.StopFollowingUser(followerID, followingID)
}

func (s *UserService) GetFollowersByUser(userID int) ([]model.UserResponse, error) {
	users, err := s.userRepo.GetFollowersByUser(userID)
	if err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, nil
}

func (s *UserService) GetFollowingByUser(userID int) ([]model.UserResponse, error) {
	users, err := s.userRepo.GetFollowingByUser(userID)
	if err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, nil
}

func (s *UserService) ProfileUpdate(userID int, updates map[string]interface{}) (*model.User, error) {
	user, err := s.userRepo.ProfileUpdate(userID, updates)
	if err != nil {
		return nil, errors.New("failed to update profile")
	}
	return user, nil
}

func (s *UserService) PasswordChange(userID int, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if !hash.CheckPassword(oldPassword, user.Password) {
		return errors.New("invalid old_password")
	}
	hashedNewPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.PasswordChange(user.UserID, hashedNewPassword)
}
