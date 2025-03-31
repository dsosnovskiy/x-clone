package service

import (
	"fmt"
	"x-clone/internal/model"
	"x-clone/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(userID int) (*model.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	userResponse := user.ToResponse()
	return &userResponse, nil
}

func (s *UserService) FindUserByUsername(username string) (*model.UserResponse, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	userResponse := user.ToResponse()
	return &userResponse, nil
}

func (s *UserService) FollowUser(followerID, followingID int) error {
	if followerID == followingID {
		return fmt.Errorf("you cannot follow or stop following yourself")
	}
	if err := s.userRepo.FollowUser(followerID, followingID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) StopFollowingUser(followerID, followingID int) error {
	if followerID == followingID {
		return fmt.Errorf("you cannot follow or stop following yourself")
	}
	if err := s.userRepo.StopFollowingUser(followerID, followingID); err != nil {
		return err
	}
	return nil
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
