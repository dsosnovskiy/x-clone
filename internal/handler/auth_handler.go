package handler

import (
	"encoding/json"
	"net/http"
	"x-clone/internal/model"
	"x-clone/internal/service"
	"x-clone/internal/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Req parsing
		var req validator.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// Validation
		if err := validator.Validate(req); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		// To model
		user := model.User{
			Username:  req.Username,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Birthday:  req.Birthday,
			Bio:       req.Bio,
		}

		// Service call
		token, err := h.authService.Register(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": token,
		})
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Req parsing
		var req validator.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// Validation
		if err := validator.Validate(req); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		// Service call
		token, err := h.authService.Login(req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": token,
		})
	}
}
