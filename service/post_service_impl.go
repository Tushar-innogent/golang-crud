// service/post_service.go
package service

import (
	"errors"
	"go-crud/models"
	"go-crud/repository"
	"log"
)

type PostServiceImpl struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) PostService {
	return &PostServiceImpl{repo: repo}
}

func (s *PostServiceImpl) CreatePost(post *models.Post) error {
	log.Println("Create post function called! post:", post)
	err := s.repo.Create(post)
	return handleRepoError(err, "failed to create post")
}

func (s *PostServiceImpl) GetPostsByUserId(userId string) ([]models.Post, error) {
	posts, err := s.repo.FindByUserId(userId)
	return posts, handleRepoError(err, "failed to get posts by user ID")
}

func (s *PostServiceImpl) GetPostById(id string) (*models.Post, error) {
	post, err := s.repo.FindById(id)
	return post, handleRepoError(err, "failed to get post by ID")
}

func handleRepoError(err error, message string) error {
	if err != nil {
		return errors.New(message + ": " + err.Error())
	}
	return nil
}
