package service

import (
	"app/internal/domain"
	"app/internal/repository"
	"errors"
)

type PostService interface {
	CreatePost(req domain.CreatePostRequest, authorID uint) (*domain.Post, error)
	GetAllPosts() ([]domain.Post, error)
	GetPostByID(id uint) (*domain.Post, error)
	UpdatePost(id uint, req domain.UpdatePostRequest, currentUserID uint) (*domain.Post, error)
	DeletePost(id uint, currentUserID uint) error
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) CreatePost(req domain.CreatePostRequest, authorID uint) (*domain.Post, error) {
	post := domain.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
	}

	err := s.postRepo.Create(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *postService) GetAllPosts() ([]domain.Post, error) {
	return s.postRepo.FindAll()
}

func (s *postService) GetPostByID(id uint) (*domain.Post, error) {
	post, err := s.postRepo.FindById(id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (s *postService) UpdatePost(id uint, req domain.UpdatePostRequest, currentUserID uint) (*domain.Post, error) {
	post, err := s.postRepo.FindById(id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.AuthorID != currentUserID {
		return nil, errors.New("you are not authorized to update this post")
	}

	post.Title = req.Title
	post.Content = req.Content

	err = s.postRepo.Update(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) DeletePost(id uint, currentUserID uint) error {
	post, err := s.postRepo.FindById(id)
	if err != nil {
		return errors.New("post not found")
	}

	if post.AuthorID != currentUserID {
		return errors.New("you are not authorized to delete this post")
	}

	return s.postRepo.Delete(id)
}
