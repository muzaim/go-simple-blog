package repository

import (
	"app/internal/domain"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *domain.Post) error
	FindAll() ([]domain.Post, error)
	FindById(id uint) (*domain.Post, error)
	Update(post *domain.Post) error
	Delete(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindAll() ([]domain.Post, error) {
	var posts []domain.Post

	err := r.db.Preload("Author").Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *postRepository) FindById(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Preload("Author").First(&post, id).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Post{}, id).Error
}
