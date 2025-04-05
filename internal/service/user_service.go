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

func (s *UserService) GetUserByID(userID int) (*model.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	userResponse := user.ToResponse()
	return &userResponse, nil
}

func (s *UserService) FindUserByUsername(username string) (*model.UserResponse, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	userResponse := user.ToResponse()
	return &userResponse, nil
}

func (s *UserService) FollowUser(followerID, followingID int) error {
	if followerID == followingID {
		return errors.New("you cannot follow or stop following yourself")
	}
	if err := s.userRepo.FollowUser(followerID, followingID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) StopFollowingUser(followerID, followingID int) error {
	if followerID == followingID {
		return errors.New("you cannot follow or stop following yourself")
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

func (s *UserService) ChangeProfile(userID int, username, firstName, lastName, birthday, bio *string) error {
	if err := s.userRepo.ChangeProfile(userID, username, firstName, lastName, birthday, bio); err != nil {
		return err
	}
	return nil
}

func (s *UserService) ChangePassword(userID int, oldPassword, newPassword, confirmPassword string) error {
	if oldPassword == newPassword {
		return errors.New("the new password cannot be equal to the old password")
	}
	if newPassword != confirmPassword {
		return errors.New("failed password confirmation")
	}
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if !hash.CheckPassword(oldPassword, user.Password) {
		return errors.New("invalid old_password")
	}

	hashedNewPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := s.userRepo.ChangePassword(userID, hashedNewPassword); err != nil {
		return err
	}
	return nil
}
