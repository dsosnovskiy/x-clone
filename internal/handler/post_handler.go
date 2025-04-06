package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"x-clone/internal/model"
	"x-clone/internal/service"
	"x-clone/internal/validator"
	"x-clone/pkg/middleware"

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
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Req parsing
		var req validator.ContentRequest
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
		post := model.Post{
			UserID:  userID,
			Content: req.Content,
		}

		// Service call
		newPost, err := h.postService.CreatePost(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPost)
	}
}

func (h *PostHandler) GetUserPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL parsing
		username := chi.URLParam(r, "username")

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		posts, err := h.postService.GetUserPosts(user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

func (h *PostHandler) GetUserPostByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		post, err := h.postService.GetUserPostByID(user.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func (h *PostHandler) UpdatePostContentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Req parsing
		var req validator.ContentRequest
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
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if user.UserID != userID {
			http.Error(w, "you are not owner of this post", http.StatusForbidden)
			return
		}
		post, err := h.postService.GetUserPostByID(userID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		updatedPost, err := h.postService.UpdatePostContentByID(userID, post.PostID, req.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedPost)
	}
}

func (h *PostHandler) DeletePostByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if user.UserID != userID {
			http.Error(w, "you are not owner of this post", http.StatusForbidden)
			return
		}
		post, err := h.postService.GetUserPostByID(userID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := h.postService.DeletePostByID(userID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "successfully deleted the post",
		})
	}
}

func (h *PostHandler) LikePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		post, err := h.postService.GetUserPostByID(user.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := h.postService.LikePost(userID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "successfully liked the post",
			"post_id": postID,
		})
	}
}

func (h *PostHandler) UnlikePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		post, err := h.postService.GetUserPostByID(user.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := h.postService.UnlikePost(userID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "successfully unliked the post",
			"post_id": postID,
		})
	}
}

func (h *PostHandler) RepostPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		post, err := h.postService.GetUserPostByID(user.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := h.postService.RepostPost(userID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "successfully reposted the post",
			"post_id": postID,
		})
	}
}

func (h *PostHandler) UndoRepostPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		post, err := h.postService.GetUserPostByID(user.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := h.postService.UndoRepostPost(userID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "successfully cancelled repost of the post",
			"post_id": postID,
		})
	}
}

func (h *PostHandler) GetUserReposts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL parsing
		username := chi.URLParam(r, "username")

		// Service call
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		posts, err := h.postService.GetUserReposts(user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

func (h *PostHandler) QuotePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// URL parsing
		username := chi.URLParam(r, "username")
		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "invalid post_id", http.StatusBadRequest)
			return
		}

		// Req parsing
		var req validator.ContentRequest
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
		user, err := h.userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		originalPost, err := h.postService.GetUserPostByID(user.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		quotePost, err := h.postService.QuotePost(userID, originalPost.PostID, req.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(quotePost)
	}
}
