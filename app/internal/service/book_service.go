package service

import (
	"app/internal/domain"
	"app/internal/repository"
)

type BookService interface {
	CreateBook(req domain.CreateNewBookRequest) error
	GetAllBooks() ([]domain.Book, error)
	GetBookByID(id uint) (*domain.Book, error)
	UpdateBook(id uint, req domain.UpdateBookRequest) (*domain.Book, error)
	DeleteBook(id uint) error
}

type bookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{bookRepo: bookRepo}
}

func (s *bookService) CreateBook(req domain.CreateNewBookRequest) error {
	book := domain.Book{
		Title:  req.Title,
		Isbn:   req.Isbn,
		Author: req.Author,
		Price:  req.Price,
		Stock:  req.Stock,
	}

	err := s.bookRepo.Create(&book)
	if err != nil {
		return err
	}
	return nil
}

func (s *bookService) GetAllBooks() ([]domain.Book, error) {
	books, err := s.bookRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *bookService) GetBookByID(id uint) (*domain.Book, error) {
	book, err := s.bookRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *bookService) UpdateBook(id uint, req domain.UpdateBookRequest) (*domain.Book, error) {
	book, err := s.bookRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	book.Title = req.Title
	book.Isbn = req.Isbn
	book.Author = req.Author
	book.Price = req.Price
	book.Stock = req.Stock

	err = s.bookRepo.Update(book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) DeleteBook(id uint) error {
	return s.bookRepo.Delete(id)
}
