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
	bookRepo := repository.NewBookRepository(db)

	authService := service.NewAuthService(userRepo)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)
	bookService := service.NewBookService(bookRepo)

	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	bookHandler := handler.NewBookHandler(bookService)

	r := router.SetupRouter(router.RouterConfig{
		AuthHandler:    authHandler,
		PostHandler:    postHandler,
		CommentHandler: commentHandler,
		BookHandler:    bookHandler,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	r.Run(":" + port)
}
