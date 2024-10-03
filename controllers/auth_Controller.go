package controllers

import (
	"go-crud/auth"
	"go-crud/models"
	"go-crud/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService service.UserService
}

func NewAuthController(userService service.UserService) *AuthController {
	return &AuthController{userService: userService}
}
func (uc *AuthController) Signup(c *gin.Context) {
	var body struct {
		Name      string
		Email     string
		Password  string
		CompanyID uint
		Role      string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//create the user
	user := models.User{Name: body.Name, CompanyID: body.CompanyID, Email: body.Email, Password: body.Password, Role: models.ParseRole(body.Role)}
	result, err := uc.userService.CreateUser(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"User registered ": result.Name})
}

func (ac *AuthController) Login(c *gin.Context) {
	var requestBody models.LoginRequest

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user *models.User
	var err error
	user, err = ac.userService.FindByEmail(requestBody.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token for authenticated user
	tokenString, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error generating token",
		})
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
