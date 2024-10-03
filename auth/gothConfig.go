package auth

import (
	"fmt"
	"go-crud/service"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type GothConfig struct {
	userService service.UserService
}

func NewGothConfig(userService service.UserService) *GothConfig {
	return &GothConfig{userService: userService}
}

func ConfigGoth() {

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	clientCallbackURL := os.Getenv("GOOGLE_CLIENT_CALLBACK_URL")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	githubClientCallbackURL := os.Getenv("GITHUB_CLIENT_CALLBACK_URL")

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
		github.New(githubClientID, githubClientSecret, githubClientCallbackURL),
	)

}

func (gc *GothConfig) SignInWithProvider(c *gin.Context) {

	provider := c.Param("provider")
	log.Println("Auth Provider : ", provider)
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (gc *GothConfig) CallbackHandler(c *gin.Context) {

	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println("User : ", user)
	userData, err := gc.userService.FindByEmail(user.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User doesn't exists with email: " + user.Email,
		})
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	tokenString, err := GenerateToken(userData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error generating token",
		})
	}

	fmt.Println("Token string : ", tokenString)

	c.Redirect(http.StatusTemporaryRedirect, "/success")
}

func (gc *GothConfig) Success(c *gin.Context) {

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`
      <div style="
          background-color: #fff;
          padding: 40px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          text-align: center;
      ">
          <h1 style="
              color: #333;
              margin-bottom: 20px;
          ">You have Successfully signed in!</h1>
          
          </div>
      </div>
  `)))
}
