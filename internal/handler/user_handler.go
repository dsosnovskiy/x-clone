package handler

import (
	"encoding/json"
	"net/http"
	"x-clone/internal/model"
	"x-clone/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := user.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.userService.CreateUser(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userResponse := user.ToResponse()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(userResponse)
	}
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
		followingUsername := chi.URLParam(r, "username")

		followerUsername := r.Header.Get("Username")

		if followerUsername == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}

		followerUser, err := h.userService.FindUserByUsername(followerUsername)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		followingUser, err := h.userService.FindUserByUsername(followingUsername)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.userService.FollowUser(followerUser.UserID, followingUser.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successfully follow user: " + followingUsername))
	}
}

func (h *UserHandler) StopFollowingUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followingUsername := chi.URLParam(r, "username")

		followerUsername := r.Header.Get("Username")

		if followerUsername == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}

		followerUser, err := h.userService.FindUserByUsername(followerUsername)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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
		w.Write([]byte("successfully stop following user: " + followingUsername))
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
