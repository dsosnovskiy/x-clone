package service

import (
	"errors"
	"time"
	"x-clone/internal/config"
	"x-clone/internal/model"
	"x-clone/internal/repository"
	"x-clone/pkg/utils/hash"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(authRepo *repository.AuthRepository, userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{authRepo: authRepo, userRepo: userRepo, cfg: cfg}
}

func (s *AuthService) Register(user *model.User) (string, error) {
	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	user.Password = hashedPassword

	if err := s.authRepo.CreateUser(user); err != nil {
		return "", errors.New("username already exists")
	}

	return s.GenerateAccessToken(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !hash.CheckPassword(password, user.Password) {
		return "", errors.New("invalid password")
	}

	return s.GenerateAccessToken(user)
}

func (s *AuthService) GenerateAccessToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.UserID,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"exp":        time.Now().Add(s.cfg.JWT.AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *AuthService) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	if s.cfg.JWT.Secret == "" {
		return nil, errors.New("JWT secret not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
