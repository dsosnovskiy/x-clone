package main

import (
	"log"
	"net/http"
	"x-clone/internal/config"
	"x-clone/internal/handler"
	"x-clone/internal/repository"
	"x-clone/internal/router"
	"x-clone/internal/service"
	"x-clone/pkg/database"
	"x-clone/pkg/logging"
	"x-clone/pkg/middleware"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

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
	authRepo := repository.NewAuthRepository(db)
	log.Debug("Successfully initialized the repository")

	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo, userRepo)
	authService := service.NewAuthService(authRepo, userRepo, cfg)
	log.Debug("Successfully initialized the service")

	authMiddleware := middleware.AuthMiddleware(authService)
	log.Debug("Successfully initialized middleware")

	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService, userService)
	authHandler := handler.NewAuthHandler(authService)
	log.Debug("Successfully initialized the handler")

	handlers := &router.Handlers{
		AuthHandler: authHandler,
		PostHandler: postHandler,
		UserHandler: userHandler,
	}
	r := router.New(handlers, authMiddleware)
	log.Debug("Successfully initialized the router")

	log.Infof("The server is running on address: %s", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, r); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
