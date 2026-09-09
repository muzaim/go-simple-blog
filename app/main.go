package main

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/repository"
	"app/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// http.HandleFunc("/", handler)
	// fmt.Println("Server is running on http://localhost:8080")
	// http.ListenAndServe(":8080", nil)

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

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPostByID)
	r.GET("/posts/:id/comments", commentHandler.GetComments)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
		protected.POST("/posts/:id/comments", commentHandler.CreateComment)
	}
	fmt.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
