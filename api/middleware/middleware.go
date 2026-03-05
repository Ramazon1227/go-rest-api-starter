package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	httpapi "github.com/Ramazon1227/go-rest-api-starter/api/http"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(httpapi.Unauthorized.Code, httpapi.Response{
				Status:    httpapi.Unauthorized.Status,
				Description: "authorization header is required",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(httpapi.Unauthorized.Code, httpapi.Response{
				Status:    httpapi.Unauthorized.Status,
				Description: "authorization header is required",
			})
			return
		}

		token := parts[1]
		if jwt.IsTokenBlacklisted(token) {
			c.AbortWithStatusJSON(httpapi.Unauthorized.Code, httpapi.Response{
				Status:    httpapi.Unauthorized.Status,
				Description: "token has been invalidated",
			})		
			return
		}

		claims, err := jwt.ExtractClaims(token)
		if err != nil {
			c.AbortWithStatusJSON(httpapi.Unauthorized.Code, httpapi.Response{
				Status:    httpapi.Unauthorized.Status,
				Description: err.Error(),
			})
			return
		}

		// Set user information in context
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Set("email", claims["email"])

		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(httpapi.Unauthorized.Code, httpapi.Response{
				Status:    httpapi.Unauthorized.Status,
				Description: "role not found in token",
			})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(httpapi.Unauthorized.Code, httpapi.Response{
				Status:      httpapi.Unauthorized.Status,
				Description: "invalid role type in token",
			})
			return
		}
		allowed := false
		for _, r := range allowedRoles {
			if r == roleStr {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(httpapi.Forbidden.Code, httpapi.Response{
				Status:    httpapi.Forbidden.Status,
				Description: "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
