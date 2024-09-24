package main

import (
	"go-crud/controllers"
	"go-crud/initializers"
	"go-crud/repository"
	"go-crud/service"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// Initialize database connection
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

	// Set up repository and services
	repo := repository.NewCompanyRepository(initializers.DB) // Assume this is correctly implemented
	companyService := service.NewCompanyServiceImpl(repo)
	companyController := controllers.NewCompanyController(companyService)

	postRepo := repository.NewPostRepository(initializers.DB)
	postService := service.NewPostService(postRepo)
	postController := controllers.NewPostController(postService)

	userRepo := repository.NewUserRepository(initializers.DB)
	userService := service.NewUserServiceImpl(userRepo) // Returns an implementation of UserService interface
	userController := controllers.NewUserController(userService)

	// Create a Gin router
	r := gin.Default()

	//Company API's
	r.POST("/company", companyController.CreateCompany)
	r.GET("/getAllCompanies", companyController.GetAllCompanies)
	r.DELETE("/deleteCompany/:id", companyController.DeleteCompany)

	//Post API's
	r.POST("/post", postController.CreatePost)
	r.GET("/getAllPosts/:id", postController.GetPosts)
	r.GET("/getPost/:id", postController.GetPostById)

	//Users API's
	r.POST("/user", userController.CreateUser)
	r.GET("/getUsers", userController.GetUsers)
	r.GET("/getUserById/:id", userController.GetUserById)
	r.PUT("/updateUser/:id", userController.UpdateUserDetails)
	r.DELETE("/deleteUser/:id", userController.DeleteUser)
	r.GET("/paginatedUser", userController.PaginateUsers)

	r.Run()
}
