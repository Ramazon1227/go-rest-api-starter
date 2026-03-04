package main

import (
	"log"
	"net/http"

	"github.com/Ramazon1227/go-rest-api-starter/internal/config"
	"github.com/Ramazon1227/go-rest-api-starter/internal/handler"
	"github.com/Ramazon1227/go-rest-api-starter/internal/repository"
	"github.com/Ramazon1227/go-rest-api-starter/internal/server"
	"github.com/Ramazon1227/go-rest-api-starter/internal/service"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/database"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.NewPool(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, jwtManager)

	authHandler := handler.NewAuthHandler(userSvc)
	userHandler := handler.NewUserHandler(userSvc)

	r := server.New(authHandler, userHandler, jwtManager)

	addr := ":" + cfg.ServerPort
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
