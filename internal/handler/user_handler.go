package handler

import (
	"encoding/json"
	"net/http"
	"x-clone/internal/model"
	"x-clone/internal/service"
	"x-clone/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) FindUserByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		userResponse, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userResponse)
	}
}

func (h *UserHandler) FollowUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		followingUsername := chi.URLParam(r, "username")

		followingUser, err := h.userService.FindUserByUsername(followingUsername)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.userService.FollowUser(userID, followingUser.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful follow user: " + followingUsername))
	}
}

func (h *UserHandler) StopFollowingUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		followingUsername := chi.URLParam(r, "username")

		followerUser, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		followingUser, err := h.userService.FindUserByUsername(followingUsername)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.userService.StopFollowingUser(followerUser.UserID, followingUser.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful stop following user: " + followingUsername))
	}
}

func (h *UserHandler) GetFollowersByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		followers, err := h.userService.GetFollowersByUser(user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if followers == nil {
			followers = []model.UserResponse{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(followers)
	}
}

func (h *UserHandler) GetFollowingByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		following, err := h.userService.GetFollowingByUser(user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if following == nil {
			following = []model.UserResponse{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(following)
	}
}

func (h *UserHandler) ChangeProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var userData struct {
			Username  *string `json:"username"`
			FirstName *string `json:"first_name"`
			LastName  *string `json:"last_name"`
			Birthday  *string `json:"birthday"`
			Bio       *string `json:"bio"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := model.ValidateChangeProfile(
			userData.Username,
			userData.FirstName,
			userData.LastName,
			userData.Birthday,
			userData.Bio); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.userService.ChangeProfile(
			userID,
			userData.Username,
			userData.FirstName,
			userData.LastName,
			userData.Birthday,
			userData.Bio); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func (h *UserHandler) ChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var passwordData struct {
			OldPassword     string `json:"old_password"`
			NewPassword     string `json:"new_password"`
			ConfirmPassword string `json:"confirm_password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&passwordData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := model.ValidateChangePassword(
			passwordData.OldPassword,
			passwordData.NewPassword,
			passwordData.ConfirmPassword); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.userService.ChangePassword(
			userID,
			passwordData.OldPassword,
			passwordData.NewPassword,
			passwordData.ConfirmPassword); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful password change"))
	}
}
