package main

import (
	"go-crud/auth"
	"go-crud/controllers"
	"go-crud/initializers"
	"go-crud/middleware"
	"go-crud/repository"
	"go-crud/service"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	// Initialize database connection
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	auth.ConfigGoth()

	// Configure Gothic with your session store
	gothic.Store = store

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

	authController := controllers.NewAuthController(userService)

	gothConfig := auth.NewGothConfig(userService)

	// Create a Gin router
	r := gin.Default()

	//Company API's
	r.POST("/company", companyController.CreateCompany)
	r.GET("/getAllCompanies", middleware.RequireAuth("RoleUser", "RoleAdmin"), companyController.GetAllCompanies)
	r.DELETE("/deleteCompany/:id", middleware.RequireAuth("RoleAdmin"), companyController.DeleteCompany)

	//Post API's
	r.POST("/post", middleware.RequireAuth("RoleUser", "RoleAdmin"), postController.CreatePost)
	r.GET("/getAllPosts/:id", middleware.RequireAuth("RoleUser", "RoleAdmin"), postController.GetPosts)
	r.GET("/getPost/:id", middleware.RequireAuth("RoleUser", "RoleAdmin"), postController.GetPostById)

	//Users API's
	// r.POST("/signup", userController.CreateUser)
	r.POST("/signup", authController.Signup)
	r.POST("login", authController.Login)

	//auth/google/callback

	r.GET("/getUsers", middleware.RequireAuth("RoleAdmin"), userController.GetUsers)
	r.GET("/getUserById/:id", middleware.RequireAuth("RoleUser", "RoleAdmin"), userController.GetUserById)
	r.PUT("/updateUser/:id", middleware.RequireAuth("RoleUser"), userController.UpdateUserDetails)
	r.DELETE("/deleteUser/:id", middleware.RequireAuth("RoleAdmin"), userController.DeleteUser)
	r.GET("/paginatedUser", middleware.RequireAuth("RoleUser", "RoleAdmin"), userController.PaginateUsers)
	r.POST("/singleTransac", middleware.RequireAuth("RoleUser", "RoleAdmin"), userController.SingleTransaction)

	r.GET("/", home)
	r.GET("/auth/:provider", gothConfig.SignInWithProvider)
	r.GET("/auth/:provider/callback", gothConfig.CallbackHandler)
	r.GET("/success", gothConfig.Success)

	// r.GET("/auth/github", auth.SignInWithProvider)
	// r.GET("/auth/github/callback", auth.CallbackHandler)

	r.Run()
}

func home(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(c.Writer, gin.H{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// func beginAuthHandler(c *gin.Context) {
// 	// Begin the OAuth2 authentication flow
// 	provider, err := goth.GetProvider("github")
// 	if err != nil {
// 		c.String(500, err.Error())
// 		return
// 	}

// 	// Redirect the user to the provider's authentication page
// 	url, err := provider.BeginAuth("state")
// 	if err != nil {
// 		c.String(500, err.Error())
// 		return
// 	}

// 	c.Redirect(302, url)
// }

// func callbackHandler(c *gin.Context) {
// 	// Get the user's information after authentication
// 	provider, err := goth.GetProvider("github")
// 	if err != nil {
// 		c.String(500, err.Error())
// 		return
// 	}

// 	// Retrieve the user from the request
// 	user, err := provider.GetUser(c.Request)
// 	if err != nil {
// 		c.String(500, err.Error())
// 		return
// 	}

// 	// Display the user's information
// 	c.JSON(200, user)
// }
