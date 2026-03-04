package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/response"
)

const UserIDKey = "user_id"

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtManager.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	id, _ := c.Get(UserIDKey)
	userID, _ := id.(int64)
	return userID
}
