package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"x-clone/internal/model"
	"x-clone/internal/service"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postService *service.PostService
	userService *service.UserService
}

func NewPostHandler(postService *service.PostService, userService *service.UserService) *PostHandler {
	return &PostHandler{postService: postService, userService: userService}
}

func (h *PostHandler) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")
		if username == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}

		var post model.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := post.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.postService.CreatePost(&post, username); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

func (h *PostHandler) GetUserPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		posts, err := h.postService.GetUserPosts(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

func (h *PostHandler) GetUserPostByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		post, err := h.postService.GetUserPostByID(user.UserID, int(postID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func (h *PostHandler) UpdatePostContentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")

		if username == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var content struct {
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := h.postService.UpdatePostContentByID(int(postID), user.UserID, content.Content); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		updatedPost, err := h.postService.GetUserPostByID(user.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedPost)
	}
}

func (h *PostHandler) DeletePostByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")

		if username == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.postService.DeletePostByID(int(postID), user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
