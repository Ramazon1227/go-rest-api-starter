package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ramazon1227/go-rest-api-starter/api"
	"github.com/Ramazon1227/go-rest-api-starter/api/handlers"
	"github.com/Ramazon1227/go-rest-api-starter/config"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/email"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/logger"
	"github.com/Ramazon1227/go-rest-api-starter/storage/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	// Wire config to packages
	jwt.SigningKey = []byte(cfg.SecretKey)

	email.Host = cfg.SMTPHost
	email.Port = fmt.Sprintf("%d", cfg.SMTPPort)
	email.Username = cfg.SMTPUsername
	email.Password = cfg.SMTPPassword
	email.From = cfg.SMTPFrom

	pgStore, err := postgres.NewPostgres(context.Background(), cfg)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	defer pgStore.CloseDB()

	h := handlers.NewHandler(cfg, log, pgStore)

	r := api.SetUpRouter(h, cfg)

	srv := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: r,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic("srv.ListenAndServe", logger.Error(err))
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Panic("srv.Shutdown", logger.Error(err))
	}
}
