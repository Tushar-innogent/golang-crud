package test

import (
	"go-crud/models"
	"go-crud/repository"
	"log"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB() (*gorm.DB, error) {
	// Configure GORM logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // IO writer (stdout)
		logger.Config{
			SlowThreshold: time.Second, // Log SQL queries that take longer than this
			LogLevel:      logger.Info, // LogLevel: Silent, Error, Warn, Info
			Colorful:      true,        // Enable/Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	return db, nil
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewUserRepository(db)
	user := &models.User{Name: "John Doe"}

	// Test Create method
	createdUser, err := repo.Create(user)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", createdUser.Name)
	assert.NotZero(t, createdUser.ID)
	db.Where("1 = 1").Delete(&models.User{})
}

func Test_FindAll(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewUserRepository(db)
	user := &models.User{Name: "John", Email: "john@gmail.com"}
	expectedUsers := []*models.User{user}

	repo.Create(user)

	result, err := repo.FindAll()

	log.Println("expectedUsers[0].Name : ", expectedUsers[0].Name)
	log.Println("result[0].Name : ", result[0].Name)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers[0].Name, result[0].Name)
	assert.NotZero(t, result)

	db.Where("1 = 1").Delete(&models.User{})
}

func TestMultipleUpdateSaveTransaction_Error(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	user := &models.User{Name: "Steve", Email: "steve@email.com"}
	db.Create(user)

	repo := repository.NewUserRepository(db)

	// Attempt the transaction which is expected to fail
	_, err = repo.MultipleUpdateSaveTransaction(user)

	// AssertionsDB:db
	assert.Error(t, err)

	// Fetch the user again to ensure no updates were applied
	var result models.User
	db.First(&result, user.ID)
	assert.Equal(t, "Steve", result.Name)            // Should remain unchanged
	assert.Equal(t, "steve@email.com", result.Email) // Should remain unchanged

	db.Where("1 = 1").Delete(&models.User{})
}

func Test_FindByID(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	expectedUser := models.User{ID: 1, Name: "Tushar", Email: "Patidar"}

	repo := repository.UserRepositoryImpl{DB: db}

	repo.Create(&expectedUser)

	result, err := repo.FindById("1")

	assert.Error(t, err, nil)
	assert.Equal(t, expectedUser.Name, result.Name)
}

func Test_Update(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	user := models.User{Name: "Tushar", Email: "tushar@gmail.com"}
	data := map[string]interface{}{"Name": "Tushar Patidar", "Email": "tushar.patidar@gmail.com"}

	repo := repository.NewUserRepository(db)
	repo.Create(&user)

	repo.Update(&user, data)

	assert.Equal(t, user.Name, data["Name"])

	db.Where("1 = 1").Delete(&models.User{})
}

func Test_Delete(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	var users []models.User
	users = append(users, models.User{ID: 1, Name: "Tushar", Email: "tushar@gmail.com"})
	users = append(users, models.User{ID: 2, Name: "Jay", Email: "jay@gmail.com"})
	db.Save(users)

	repo := repository.NewUserRepository(db)
	dbError := repo.Delete("1")

	resultOfAllUsers, _ := repo.FindAll()
	log.Println("all user", resultOfAllUsers)

	assert.Equal(t, len(resultOfAllUsers), 1)
	assert.NoError(t, dbError)

	db.Where("1 = 1").Delete(&models.User{})
}

func Test_Paginate(t *testing.T) {
	offset := 0
	limit := 2

	db, err := setupTestDB()
	assert.NoError(t, err)

	var users []models.User
	users = append(users, models.User{ID: 1, Name: "Tushar", Email: "tushar@gmail.com"})
	users = append(users, models.User{ID: 2, Name: "Jay", Email: "jay@gmail.com"})
	db.Save(users)

	repo := repository.NewUserRepository(db)
	result, err := repo.Paginate(offset, limit)

	assert.Equal(t, len(users), len(result))
	assert.NoError(t, err)

	db.Where("1 = 1").Delete(&models.User{})
}
