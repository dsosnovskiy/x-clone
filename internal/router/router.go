package router

import (
	"x-clone/internal/handler"

	"github.com/go-chi/chi/v5"
)

func New(postHandler *handler.PostHandler, userHandler *handler.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/users", userHandler.CreateUser())
	r.Get("/users", userHandler.FindUserByUsername())

	r.Post("/posts", postHandler.CreatePost())
	r.Get("/posts", postHandler.GetUserPosts())

	return r
}
