package domain

import "time"

type Comment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PostID     uint      `gorm:"not null" json:"post_id"`
	AuthorName string    `gorm:"size:100;not null" json:"author_name"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateCommentRequest struct {
	AuthorName string `json:"author_name" binding:"required"`
	Content    string `json:"content" binding:"required"`
}
