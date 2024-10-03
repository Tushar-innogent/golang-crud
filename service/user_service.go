// service/user_service.go
package service

import "go-crud/models"

type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserById(id string) (*models.User, error)
	UpdateUserDetails(user *models.User, data map[string]interface{}) error
	DeleteUser(id string) error
	PaginateUsers(page, pageSize int) ([]models.User, error)
	SingleTransactionUser(user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}
