package auth

import (
	"github.com/gin-gonic/gin"
	"legato_server/api"
	"legato_server/authenticate"
	"legato_server/internal/legato/api/rest/server"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
	"legato_server/middleware"
	"net/http"
)

type Auth struct {
	db database.Database
}

func (a *Auth) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/auth/signup", a.Signup)
	group.POST("/auth/login", a.Login)
	group.GET("/auth/user", a.LoggedInUser)
}

func NewAuthModule(db database.Database) (server.RestModule, error) {
	return &Auth{
		db: db,
	}, nil
}

func (a *Auth) Signup(c *gin.Context) {
	newUser := SignupRequest{}
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create the user
	if err = a.db.AddUser(models.User{
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully.",
	})
}

func (a *Auth) Login(c *gin.Context) {
	cred := LoginRequest{}
	err := c.ShouldBindJSON(&cred)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check username validation
	expectedUser, err := a.db.GetUserByUsername(cred.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check credentials
	token, err := authenticate.Login(cred, expectedUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.TokenString,
	})
}

func (a *Auth) LoggedInUser(c *gin.Context) {
	loggedInUser := checkAuth(c, nil)
	if loggedInUser == nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": api.UserInfo{
			ID:       loggedInUser.ID,
			Email:    loggedInUser.Email,
			Username: loggedInUser.Username,
		},
	})
}

// checkAuth was written because of DRY (Don't Repeat Yourself).
// Each time it authenticate the user and handle the errors that might occur.
// validUsernames is the list of usernames that the api is accessible for them.
// nil validUsers means that any authenticated user can use api.
// Return the logged-in user.
func checkAuth(c *gin.Context, validUsernames []string) *models.User {
	// Get the user
	rawData := c.MustGet(middleware.UserKey)
	if rawData == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return nil
	}

	loginUser := rawData.(*models.User)
	if loginUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return nil
	}

	// Check if validUsernames is nil
	// If it is nil it means any authenticated user is accepted.
	if validUsernames == nil {
		return loginUser
	}

	// If it isn't, Check if the user has access.
	for _, validUser := range validUsernames {
		if loginUser.Username == validUser {
			return loginUser
		}
	}

	// If the api is not accessible
	c.JSON(http.StatusForbidden, gin.H{
		"message": "access denied: can not do any actions for this user",
	})
	return nil
}
