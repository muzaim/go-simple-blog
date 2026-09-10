package main

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/router"
	"app/internal/service"
	"fmt"
	"os"
)

func main() {
	db := config.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	authService := service.NewAuthService(userRepo)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	r := router.SetupRouter(router.RouterConfig{
		AuthHandler:    authHandler,
		PostHandler:    postHandler,
		CommentHandler: commentHandler,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	r.Run(":" + port)
}
