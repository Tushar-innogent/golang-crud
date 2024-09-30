package test

import (
	"bytes"
	"encoding/json"
	"go-crud/controllers"
	"go-crud/models"
	"go-crud/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupPostController() (*gin.Engine, *controllers.PostController, *service.MockPostService) {
	mockPostService := new(service.MockPostService)
	postController := controllers.NewPostController(mockPostService)
	router := gin.Default()

	router.POST("/post", postController.CreatePost)
	router.GET("/getAllPosts/:id", postController.GetPosts)
	router.GET("/getPost/:id", postController.GetPostById)

	return router, postController, mockPostService
}

func Test_CreatePost(t *testing.T) {
	router, _, mockPostService := setupPostController()

	mockPost := &models.Post{Title: "Post 1", Body: "post 1 body", UserId: 1}
	mockPostService.On("CreatePost", mockPost).Return(nil) 

	postPayload, _ := json.Marshal(mockPost)
	req, _ := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(postPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	mockPostService.AssertExpectations(t)
}

func Test_GetPosts(t *testing.T){
	router, _, mockPostService := setupPostController()

	var id uint = 1
	mockPost := &models.Post{Title: "Post 1", Body: "post 1 body", UserId: 1}
	expectedResult := []models.Post{*mockPost}

	mockPostService.On("GetPosts", mockPost).Return(expectedResult)

	postPayload, _ := json.Marshal(id)
	req, _ := http.NewRequest(http.MethodPost, "/getAllPosts/:1", bytes.NewBuffer(postPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusCreated, w.Code)
	// mockPostService.AssertExpectations(t)

}