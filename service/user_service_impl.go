package service

import (
	"errors"
	"fmt"
	"go-crud/custom_error"
	"go-crud/models"
	"go-crud/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository // Keep the concrete type
}

func NewUserServiceImpl(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	result, err := s.repo.Create(user)
	if err != nil {
		return nil, err // Ensure this line exists
	}
	return result, nil
}

func (s *UserServiceImpl) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all users: %w", err)
	}
	return users, nil
}

func (s *UserServiceImpl) GetUserById(id string) (*models.User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		if errors.Is(err, custom_error.ErrUserNotFound) {
			return nil, fmt.Errorf("user with ID %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to retrieve user with ID %s: %w", id, err)
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateUserDetails(user *models.User, data map[string]interface{}) error {
	err := s.repo.Update(user, data)
	if err != nil {
		if errors.Is(err, custom_error.ErrUserNotFound) {
			return fmt.Errorf("user with ID %v not found: %w", user.ID, err)
		}
		return fmt.Errorf("failed to update user details: %w", err)
	}
	return nil
}

func (s *UserServiceImpl) DeleteUser(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, custom_error.ErrUserNotFound) {
			return fmt.Errorf("user with ID %s not found: %w", id, err)
		}
		return fmt.Errorf("failed to delete user with ID %s: %w", id, err)
	}
	return nil
}

func (s *UserServiceImpl) PaginateUsers(page, pageSize int) ([]models.User, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("invalid page or page size")
	}

	offset := (page - 1) * pageSize
	users, err := s.repo.Paginate(offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to paginate users: %w", err)
	}
	return users, nil
}
