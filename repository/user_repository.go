package repository

import "go-crud/models"

// UserRepository defines the methods for user repository operations.
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindAll() ([]models.User, error)
	FindById(id string) (*models.User, error)
	Update(user *models.User, data map[string]interface{}) error
	Delete(id string) error

	// Takes an offset and limit for pagination, and returns a slice of users or an error.
	Paginate(offset, limit int) ([]models.User, error) 
	
	// MultipleUpdateSaveTransaction updates multiple fields of a user record
    // and commits the changes within a single transaction.
	MultipleUpdateSaveTransaction(user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}
