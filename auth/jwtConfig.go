package auth

import (
	"go-crud/models"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user *models.User) (string, error) {
	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(), // Token expiration time (30 days)
		"role": models.Role.String(user.Role),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println("Error generating token:", err)
		return "", err
	}

	return tokenString, nil
}
