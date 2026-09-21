package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AuthHandler    *handler.AuthHandler
	PostHandler    *handler.PostHandler
	CommentHandler *handler.CommentHandler
	BookHandler    *handler.BookHandler
}

func SetupRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	r.POST("/register", cfg.AuthHandler.Register)
	r.POST("/login", cfg.AuthHandler.Login)
	r.GET("/posts", cfg.PostHandler.GetPosts)
	r.GET("/posts/:id", cfg.PostHandler.GetPostByID)
	r.GET("/posts/:id/comments", cfg.CommentHandler.GetComments)

	r.POST("/books", cfg.BookHandler.CreateBook)
	r.GET("/books", cfg.BookHandler.GetBooks)
	r.GET("/books/:id", cfg.BookHandler.GetBookByID)
	r.PUT("/books/:id", cfg.BookHandler.UpdateBook)
	r.DELETE("/books/:id", cfg.BookHandler.DeleteBook)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Blog Posts
		protected.POST("/posts", cfg.PostHandler.CreatePost)
		protected.PUT("/posts/:id", cfg.PostHandler.UpdatePost)
		protected.DELETE("/posts/:id", cfg.PostHandler.DeletePost)
		// Comments
		protected.POST("/posts/:id/comments", cfg.CommentHandler.CreateComment)
	}
	return r
}
