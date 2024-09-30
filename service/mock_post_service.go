package service

import (
	"go-crud/models"

	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService interface
type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) CreatePost(post *models.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostService) GetPostsByUserId(userId string) ([]models.Post, error) {
	args := m.Called(userId)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *MockPostService) GetPostById(id string) (*models.Post, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Post), args.Error(1)
}
