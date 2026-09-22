package service

import (
	"app/internal/domain"
	"app/internal/repository"
	"errors"
)

type CommentService interface {
	CreateComment(postID uint, req domain.CreateCommentRequest, userID uint) (*domain.Comment, error)
	GetCommentsByPostID(postID uint) ([]domain.Comment, error)
}

type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *commentService) CreateComment(postID uint, req domain.CreateCommentRequest, userID uint) (*domain.Comment, error) {
	_, err := s.postRepo.FindById(postID)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	comment := domain.Comment{
		PostID:  postID,
		Content: req.Content,
		UserID:  userID,
	}

	err = s.commentRepo.Create(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (s *commentService) GetCommentsByPostID(postID uint) ([]domain.Comment, error) {
	_, err := s.postRepo.FindById(postID)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	return s.commentRepo.FindByPostID(postID)
}
