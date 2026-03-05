package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/utils"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
	"github.com/Ramazon1227/go-rest-api-starter/tests/mock"
)

func TestLogin_Success(t *testing.T) {
	hashedPw, _ := utils.HashPassword("password123")
	store := mock.NewStore()
	store.UserRepo.GetByEmailFn = func(_ context.Context, email string) (*models.User, error) {
		return &models.User{
			Id:       "user-1",
			Email:    email,
			Password: hashedPw,
			Role:     "SYSTEM_ADMIN",
		}, nil
	}

	body, _ := json.Marshal(map[string]string{
		"email":    "admin@example.com",
		"password": "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, data["token"])
}

func TestLogin_UserNotFound(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.GetByEmailFn = func(_ context.Context, _ string) (*models.User, error) {
		return nil, storage.ErrorNotFound
	}

	body, _ := json.Marshal(map[string]string{
		"email":    "nobody@example.com",
		"password": "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPw, _ := utils.HashPassword("correct-password")
	store := mock.NewStore()
	store.UserRepo.GetByEmailFn = func(_ context.Context, email string) (*models.User, error) {
		return &models.User{Id: "user-1", Email: email, Password: hashedPw, Role: "STUDENT"}, nil
	}

	body, _ := json.Marshal(map[string]string{
		"email":    "user@example.com",
		"password": "wrong-password",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")

	w := newRecorder()
	newTestRouter(mock.NewStore()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_MissingPassword(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"email": "user@example.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := newRecorder()
	newTestRouter(mock.NewStore()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogout_Success(t *testing.T) {
	user := &models.User{Id: "user-1", Role: "STUDENT", Email: "u@example.com"}
	rawToken := bearerToken(t, user)[len("Bearer "):]

	body, _ := json.Marshal(map[string]string{"token": rawToken})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := newRecorder()
	newTestRouter(mock.NewStore()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
