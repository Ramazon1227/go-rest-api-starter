package server

import (
	"github.com/gin-gonic/gin"

	"github.com/Ramazon1227/go-rest-api-starter/internal/handler"
	"github.com/Ramazon1227/go-rest-api-starter/internal/middleware"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
)

func New(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtManager *jwt.Manager) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := v1.Group("/")
		protected.Use(middleware.Auth(jwtManager))
		{
			protected.GET("/users/me", userHandler.GetProfile)
			protected.PUT("/users/me", userHandler.UpdateProfile)
		}
	}

	return r
}
