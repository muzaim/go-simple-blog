package repository

import (
	"app/internal/domain"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *domain.Comment) error
	FindByPostID(postID uint) ([]domain.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return err
	}
	return r.db.Preload("User").First(comment, comment.ID).Error
}

func (r *commentRepository) FindByPostID(postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.Preload("User").Where("post_id=?", postID).Order("created_at desc").Find(&comments).Error
	return comments, err
}
