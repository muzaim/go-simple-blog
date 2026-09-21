package repository

import (
	"app/internal/domain"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *domain.Book) error
	FindAll() ([]domain.Book, error)
	FindById(id uint) (*domain.Book, error)
	Update(book *domain.Book) error
	Delete(id uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *domain.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) FindAll() ([]domain.Book, error) {
	var books []domain.Book

	err := r.db.Find(&books).Error

	return books, err
}

func (r *bookRepository) FindById(id uint) (*domain.Book, error) {
	var book domain.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}

	return &book, err
}

func (r *bookRepository) Update(book *domain.Book) error {
	return r.db.Save(&book).Error
}

func (r *bookRepository) Delete(bookId uint) error {
	return r.db.Delete(&domain.Book{}, bookId).Error
}
