package middleware

import (
	"legato_server/authenticate"
	"legato_server/internal/legato/database"
	"legato_server/pkg/logger"

	"github.com/gin-gonic/gin"
)

const Authorization = "Authorization"
const UserKey = "UserKey"

var log, _ = logger.NewLogger(logger.Config{})

type AuthMiddleware struct {
	db database.Database
}

func NewAuthMiddleware(db database.Database) GinMiddleware {
	return &AuthMiddleware{db: db}
}

func (am *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(Authorization)

		// Allow unauthenticated users in
		if token == "" {
			c.Set(UserKey, nil)
			c.Next()
			return
		}

		// Check validation jwt token
		claim, err := authenticate.CheckToken(token)
		if err != nil {
			c.Set(UserKey, nil)
			c.Next()
			return
		}

		// Get user and check if the user exists in postgres
		user, err := am.db.GetUserByUsername(claim.Username)
		if err != nil {
			c.Set(UserKey, nil)
			c.Next()
			return
		}

		// put user in context
		c.Set(UserKey, &user)
		c.Next()
	}
}
