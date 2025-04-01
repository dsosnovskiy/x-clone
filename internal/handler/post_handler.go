package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"x-clone/internal/model"
	"x-clone/internal/service"
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
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
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

		if err := h.postService.CreatePost(&post, user.UserID); err != nil {
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

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		posts, err := h.postService.GetUserPosts(user.UserID)
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
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		post, err := h.postService.GetUserPostByID(user.UserID, postID)
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
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		usernameOwner := chi.URLParam(r, "username")
		userOwner, err := h.userService.FindUserByUsername(usernameOwner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if usernameOwner != user.Username {
			http.Error(w, "you are not owner of this post", http.StatusForbidden)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		post, err := h.postService.GetUserPostByID(userOwner.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var content struct {
			Content string `json:"content" validate:"required,min=1,max=1000"`
		}

		if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := model.ValidateContent(content.Content); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.postService.UpdatePostContentByID(userOwner.UserID, post.PostID, content.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		updatedPost, err := h.postService.GetUserPostByID(userOwner.UserID, post.PostID)
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
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		usernameOwner := chi.URLParam(r, "username")
		userOwner, err := h.userService.FindUserByUsername(usernameOwner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if usernameOwner != user.Username {
			http.Error(w, "you are not owner of this post", http.StatusForbidden)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		post, err := h.postService.GetUserPostByID(userOwner.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.postService.DeletePostByID(userOwner.UserID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful deletion of the post"))
	}
}

func (h *PostHandler) LikePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userWhoLiked, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		usernameWhosePost := chi.URLParam(r, "username")
		userWhosePost, err := h.userService.FindUserByUsername(usernameWhosePost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		post, err := h.postService.GetUserPostByID(userWhosePost.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.postService.LikePost(userWhoLiked.UserID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful liking of the post: " + chi.URLParam(r, "post_id")))
	}
}

func (h *PostHandler) UnlikePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userWhoLiked, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		usernameWhosePost := chi.URLParam(r, "username")
		userWhosePost, err := h.userService.FindUserByUsername(usernameWhosePost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		post, err := h.postService.GetUserPostByID(userWhosePost.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.postService.UnlikePost(userWhoLiked.UserID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful unliking of the post: " + chi.URLParam(r, "post_id")))
	}
}

func (h *PostHandler) RepostPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userWhoReposted, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		usernameWhosePost := chi.URLParam(r, "username")
		userWhosePost, err := h.userService.FindUserByUsername(usernameWhosePost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		post, err := h.postService.GetUserPostByID(userWhosePost.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.postService.RepostPost(userWhoReposted.UserID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful repost of the post: " + chi.URLParam(r, "post_id")))
	}
}

func (h *PostHandler) UndoRepostPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userWhoReposted, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		usernameWhosePost := chi.URLParam(r, "username")
		userWhosePost, err := h.userService.FindUserByUsername(usernameWhosePost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		post, err := h.postService.GetUserPostByID(userWhosePost.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := h.postService.UndoRepostPost(userWhoReposted.UserID, post.PostID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful undorepost of the post: " + chi.URLParam(r, "post_id")))
	}
}

func (h *PostHandler) GetUserReposts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		user, err := h.userService.FindUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		posts, err := h.postService.GetUserReposts(user.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

func (h *PostHandler) QuotePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userWhoQuoted, err := h.userService.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		usernameWhosePost := chi.URLParam(r, "username")
		userWhosePost, err := h.userService.FindUserByUsername(usernameWhosePost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}

		var request struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := model.ValidateContent(request.Content); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		originalPost, err := h.postService.GetUserPostByID(userWhosePost.UserID, postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		quotePost, err := h.postService.QuotePost(userWhoQuoted.UserID, originalPost.PostID, request.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(quotePost)
	}
}
