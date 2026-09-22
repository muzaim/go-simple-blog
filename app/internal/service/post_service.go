package service

import (
	"app/internal/domain"
	"app/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type PostService interface {
	CreatePost(req domain.CreatePostRequest, authorID uint) (*domain.Post, error)
	GetAllPosts() ([]domain.Post, error)
	GetPostByID(id uint) (*domain.Post, error)
	UpdatePost(id uint, req domain.UpdatePostRequest, currentUserID uint) (*domain.Post, error)
	DeletePost(id uint, currentUserID uint) error
	BulkCreatePosts(req domain.BulkCreatePostRequest, authorID uint) (domain.BulkCreatePostResponse, error)
	BulkCreatePostsSequential(req domain.BulkCreatePostRequest, authorID uint) (domain.BulkCreatePostResponse, error)
}

type postService struct {
	postRepo repository.PostRepository
	rdb      *redis.Client
}

type postJob struct {
	Index int
	Req   domain.CreatePostRequest
}

func NewPostService(postRepo repository.PostRepository, rdb *redis.Client) PostService {
	return &postService{
		postRepo: postRepo,
		rdb:      rdb,
	}
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

	ctx := context.Background()
	s.rdb.Del(ctx, "posts:all")

	return &post, nil
}

func (s *postService) GetAllPosts() ([]domain.Post, error) {
	ctx := context.Background()
	cacheKey := "posts:all"

	cachedPosts, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []domain.Post
		if err := json.Unmarshal([]byte(cachedPosts), &posts); err == nil {
			fmt.Println("⚡ [REDIS] DATA DIAMBIL DARI CACHE REDIS!")
			return posts, nil
		}
	}

	fmt.Println("🐢 [DATABASE] DATA DIAMBIL DARI MYSQL DATABASE!")

	posts, err := s.postRepo.FindAll()
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(posts)
	if err == nil {
		s.rdb.Set(ctx, cacheKey, jsonBytes, 10*time.Minute)
	}

	return posts, nil
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

	ctx := context.Background()
	s.rdb.Del(ctx, "posts:all")

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

	err = s.postRepo.Delete(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	s.rdb.Del(ctx, "posts:all")

	return nil
}

func (s *postService) postWorker(
	id int,
	authorID uint,
	jobs <-chan postJob,
	results chan<- domain.PostJobResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for job := range jobs {
		post := domain.Post{
			Title:    job.Req.Title,
			Content:  job.Req.Content,
			AuthorID: authorID,
		}
		// Simpan ke DB
		err := s.postRepo.Create(&post)
		// Kirim hasil pemrosesan ke results channel
		results <- domain.PostJobResult{
			Post:  &post,
			Err:   err,
			Index: job.Index,
		}
	}
}

func (s *postService) BulkCreatePosts(req domain.BulkCreatePostRequest, authorID uint) (domain.BulkCreatePostResponse, error) {
	totalJobs := len(req.Posts)
	if totalJobs == 0 {
		return domain.BulkCreatePostResponse{
			TotalProcessed: 0,
			TotalSuccess:   0,
			TotalFailed:    0,
		}, nil
	}

	numWorkers := 5
	if totalJobs < numWorkers {
		numWorkers = totalJobs
	}

	jobs := make(chan postJob, totalJobs)
	results := make(chan domain.PostJobResult, totalJobs)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go s.postWorker(w, authorID, jobs, results, &wg)
	}

	for i, postReq := range req.Posts {
		jobs <- postJob{
			Index: i,
			Req:   postReq,
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var successCount, failedCount int
	for res := range results {
		if res.Err != nil {
			failedCount++
		} else {
			successCount++
		}
	}

	ctx := context.Background()
	s.rdb.Del(ctx, "posts:all")

	return domain.BulkCreatePostResponse{
		TotalProcessed: totalJobs,
		TotalSuccess:   successCount,
		TotalFailed:    failedCount,
	}, nil
}

func (s *postService) BulkCreatePostsSequential(req domain.BulkCreatePostRequest, authorID uint) (domain.BulkCreatePostResponse, error) {
	totalJobs := len(req.Posts)
	var successCount, failedCount int

	for _, postReq := range req.Posts {
		post := domain.Post{
			Title:    postReq.Title,
			Content:  postReq.Content,
			AuthorID: authorID,
		}
		if err := s.postRepo.Create(&post); err != nil {
			failedCount++
		} else {
			successCount++
		}
	}

	ctx := context.Background()
	s.rdb.Del(ctx, "posts:all")

	return domain.BulkCreatePostResponse{
		TotalProcessed: totalJobs,
		TotalSuccess:   successCount,
		TotalFailed:    failedCount,
	}, nil
}


