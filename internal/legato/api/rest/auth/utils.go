package auth

import (
	"github.com/gin-gonic/gin"
	"legato_server/internal/legato/database/models"
	"legato_server/internal/middleware"
	"net/http"
)

// CheckAuth was written because of DRY (Don't Repeat Yourself).
// Each time it authenticate the user and handle the errors that might occur.
// validUsernames is the list of usernames that the api is accessible for them.
// nil validUsers means that any authenticated user can use api.
// Return the logged-in user.
func CheckAuth(c *gin.Context, validUsernames []string) *models.User {
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
