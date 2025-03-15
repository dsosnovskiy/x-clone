package main

import (
	"net/http"
	"x-clone/internal/config"
	"x-clone/internal/handler"
	"x-clone/internal/repository"
	"x-clone/internal/router"
	"x-clone/internal/service"
	"x-clone/pkg/database"
	"x-clone/pkg/logging"
)

func main() {
	cfg := config.Load()

	log := logging.Init(cfg.Env)
	log.WithField("env", cfg.Env).Info("Starting X-clone...")

	db, err := database.ConnectDB(cfg)
	if err != nil {

		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Debug("Successfully connected to the database")

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	log.Debug("Successfully initialized the repository")

	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo, userRepo)
	log.Debug("Successfully initialized the service")

	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService, userService)
	log.Debug("Successfully initialized the handler")

	r := router.New(postHandler, userHandler)
	log.Debug("Successfully initialized the router")

	log.Infof("The server is running on address: %s", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, r); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
