// service/post_service.go
package service

import (
	"go-crud/models"
)

type PostService interface {
	CreatePost(post *models.Post) error
	GetPostsByUserId(userId string) ([]models.Post, error)
	GetPostById(id string) (*models.Post, error)
}
