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
		usernameHeader := r.Header.Get("Username")
		if usernameHeader == "" {
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

		if err := h.postService.CreatePost(&post, usernameHeader); err != nil {
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
		usernameHeader := r.Header.Get("Username")
		if usernameHeader == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}

		usernameOwner := chi.URLParam(r, "username")
		userOwner, err := h.userService.FindUserByUsername(usernameOwner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if usernameOwner != usernameHeader {
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
		usernameHeader := r.Header.Get("Username")
		if usernameHeader == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}

		usernameOwner := chi.URLParam(r, "username")
		userOwner, err := h.userService.FindUserByUsername(usernameOwner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if usernameOwner != usernameHeader {
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
	}
}

func (h *PostHandler) LikePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usernameWhoLiked := r.Header.Get("Username")
		if usernameWhoLiked == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}
		userWhoLiked, err := h.userService.FindUserByUsername(usernameWhoLiked)
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
	}
}

func (h *PostHandler) UnlikePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usernameWhoLiked := r.Header.Get("Username")
		if usernameWhoLiked == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}
		userWhoLiked, err := h.userService.FindUserByUsername(usernameWhoLiked)
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
	}
}

func (h *PostHandler) RepostPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usernameWhoReposted := r.Header.Get("Username")
		if usernameWhoReposted == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}
		userWhoReposted, err := h.userService.FindUserByUsername(usernameWhoReposted)
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
	}
}

func (h *PostHandler) UndoRepostPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usernameWhoReposted := r.Header.Get("Username")
		if usernameWhoReposted == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}
		userWhoReposted, err := h.userService.FindUserByUsername(usernameWhoReposted)
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
		usernameWhoQuoted := r.Header.Get("Username")
		if usernameWhoQuoted == "" {
			http.Error(w, "missing Username header", http.StatusUnauthorized)
			return
		}
		userWhoQuoted, err := h.userService.FindUserByUsername(usernameWhoQuoted)
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
