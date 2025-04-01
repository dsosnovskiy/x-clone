package router

import (
	"net/http"
	"x-clone/internal/handler"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	AuthHandler *handler.AuthHandler
	PostHandler *handler.PostHandler
	UserHandler *handler.UserHandler
}

func New(handlers *Handlers, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Access only
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", handlers.AuthHandler.Register())
		r.Post("/auth/login", handlers.AuthHandler.Login())
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware) // Apply middleware to all routers in the group

		// POST
		r.Post("/compose/post", handlers.PostHandler.CreatePost())
		r.Get("/{username}/posts", handlers.PostHandler.GetUserPosts())
		r.Get("/{username}/posts/{post_id}", handlers.PostHandler.GetUserPostByID())
		r.Patch("/{username}/posts/{post_id}", handlers.PostHandler.UpdatePostContentByID())
		r.Delete("/{username}/posts/{post_id}", handlers.PostHandler.DeletePostByID())
		r.Get("/{username}/reposts", handlers.PostHandler.GetUserReposts())
		r.Post("/{username}/posts/{post_id}/like", handlers.PostHandler.LikePost())
		r.Delete("/{username}/posts/{post_id}/like", handlers.PostHandler.UnlikePost())
		r.Post("/{username}/posts/{post_id}/repost", handlers.PostHandler.RepostPost())
		r.Delete("/{username}/posts/{post_id}/repost", handlers.PostHandler.UndoRepostPost())
		r.Post("/{username}/posts/{post_id}/quote", handlers.PostHandler.QuotePost())

		// USER
		r.Get("/{username}", handlers.UserHandler.FindUserByUsername())
		r.Post("/{username}/follow", handlers.UserHandler.FollowUser())
		r.Delete("/{username}/follow", handlers.UserHandler.StopFollowingUser())
		r.Get("/{username}/followers", handlers.UserHandler.GetFollowersByUser())
		r.Get("/{username}/following", handlers.UserHandler.GetFollowingByUser())
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page not found"))
	})

	return r
}
