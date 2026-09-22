package domain

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	AuthorID  uint      `gorm:"not null" json:"author_id"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type BulkCreatePostRequest struct {
	Posts []CreatePostRequest `json:"posts" binding:"required,dive"`
}

type PostJobResult struct {
	Post  *Post
	Err   error
	Index int
}

type BulkCreatePostResponse struct {
	TotalProcessed int `json:"total_processed"`
	TotalSuccess   int `json:"total_success"`
	TotalFailed    int `json:"total_failed"`
}
