package auth

import (
	"github.com/gin-gonic/gin"
	"legato_server/api"
	"legato_server/authenticate"
	"legato_server/internal/legato/api/rest/server"
	"legato_server/internal/legato/database"
	"legato_server/internal/legato/database/models"
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
	token, err := authenticate.Login(
		authenticate.LoginCredentials{
			Username: cred.Username,
			Password: cred.Password,
		},
		expectedUser,
	)
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
	loggedInUser := CheckAuth(c, nil)
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

func NewAuthModule(db database.Database) (server.RestModule, error) {
	return &Auth{
		db: db,
	}, nil
}
