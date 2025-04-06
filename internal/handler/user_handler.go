package handler

import (
	"encoding/json"
	"net/http"
	"x-clone/internal/service"
	"x-clone/internal/validator"
	"x-clone/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL parsing
		username := chi.URLParam(r, "username")

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func (h *UserHandler) FollowUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := h.userService.FollowUser(userID, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "successfully followed the user",
			"username": username,
		})
	}
}

func (h *UserHandler) StopFollowingUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := h.userService.StopFollowingUser(userID, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "successfully stop following the user",
			"username": username,
		})
	}
}

func (h *UserHandler) GetFollowersByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL parsing
		username := chi.URLParam(r, "username")

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		followers, err := h.userService.GetFollowersByUser(user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(followers)
	}
}

func (h *UserHandler) GetFollowingByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL parsing
		username := chi.URLParam(r, "username")

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		following, err := h.userService.GetFollowingByUser(user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(following)
	}
}

func (h *UserHandler) ProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Req parsing
		var req validator.ProfileUpdateRequest
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
		updates := make(map[string]interface{})
		if req.Username != nil {
			updates["username"] = *req.Username
		}
		if req.FirstName != nil {
			updates["first_name"] = *req.FirstName
		}
		if req.LastName != nil {
			updates["last_name"] = *req.LastName
		}
		if req.Birthday != nil {
			updates["birthday"] = *req.Birthday
		}
		if req.Bio != nil {
			updates["bio"] = *req.Bio
		}

		// Service call
		user, err := h.userService.ProfileUpdate(userID, updates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func (h *UserHandler) PasswordChange() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Req parsing
		var req validator.PasswordChangeRequest
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
		if err := h.userService.PasswordChange(userID, req.OldPassword, req.NewPassword); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "successfully changed password",
		})
	}
}
