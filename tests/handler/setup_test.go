package handler_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Ramazon1227/go-rest-api-starter/api"
	"github.com/Ramazon1227/go-rest-api-starter/api/handlers"
	"github.com/Ramazon1227/go-rest-api-starter/config"
	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/logger"
	"github.com/Ramazon1227/go-rest-api-starter/tests/mock"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	jwt.SigningKey = []byte("handler-test-secret")
	m.Run()
}

// noopLogger satisfies logger.LoggerI and discards all output.
type noopLogger struct{}

func (noopLogger) Debug(msg string, fields ...logger.Field)  {}
func (noopLogger) Info(msg string, fields ...logger.Field)   {}
func (noopLogger) Warn(msg string, fields ...logger.Field)   {}
func (noopLogger) Error(msg string, fields ...logger.Field)  {}
func (noopLogger) DPanic(msg string, fields ...logger.Field) {}
func (noopLogger) Panic(msg string, fields ...logger.Field)  {}
func (noopLogger) Fatal(msg string, fields ...logger.Field)  {}

func newTestRouter(store *mock.Store) *gin.Engine {
	cfg := config.Config{
		DefaultOffset: "0",
		DefaultLimit:  "10",
	}
	h := handlers.NewHandler(cfg, noopLogger{}, store)
	return api.SetUpRouter(h, cfg)
}

func bearerToken(t *testing.T, user *models.User) string {
	t.Helper()
	token, err := jwt.GenerateToken(user)
	if err != nil {
		t.Fatalf("bearerToken: %v", err)
	}
	return "Bearer " + token
}

func newRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
