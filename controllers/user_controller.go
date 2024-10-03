// controllers/user_controller.go
package controllers

import (
	"go-crud/models"
	"go-crud/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUsers - Calls GetAllUsers method in the service
func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUserById - Calls the GetUserById method in the service
func (uc *UserController) GetUserById(c *gin.Context) {
	id := c.Param("id")

	user, err := uc.userService.GetUserById(id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserDetails - Calls the UpdateUserDetails method in the service
func (uc *UserController) UpdateUserDetails(c *gin.Context) {
	var userRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	userId := c.Param("id")

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	data := map[string]interface{}{
		"name":  userRequest.Name,
		"email": userRequest.Email,
	}

	if err := uc.userService.UpdateUserDetails(user, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

// DeleteUser - Calls the DeleteUser method in the service
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := uc.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// PaginateUsers - Calls the PaginateUsers method in the service
func (uc *UserController) PaginateUsers(c *gin.Context) {
	var requestBody struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Page <= 0 {
		requestBody.Page = 1
	}

	if requestBody.PageSize <= 0 {
		requestBody.PageSize = 10
	}

	users, err := uc.userService.PaginateUsers(requestBody.Page, requestBody.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// CreateUser - Calls the CreateUser method in the service
func (uc *UserController) SingleTransaction(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := uc.userService.SingleTransactionUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": createdUser})
}
