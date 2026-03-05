package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
	"github.com/Ramazon1227/go-rest-api-starter/tests/mock"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/utils"
)

func TestGetProfile_Success(t *testing.T) {
	user := &models.User{Id: "user-1", Role: "STUDENT", Email: "me@example.com"}
	store := mock.NewStore()
	store.UserRepo.GetByIdFn = func(_ context.Context, pKey *models.PrimaryKey) (*models.User, error) {
		return &models.User{Id: pKey.Id, Name: "Me", Email: "me@example.com"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	req.Header.Set("Authorization", bearerToken(t, user))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProfile_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	// No Authorization header

	w := newRecorder()
	newTestRouter(mock.NewStore()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateProfile_Success(t *testing.T) {
	user := &models.User{Id: "user-1", Role: "STUDENT", Email: "me@example.com"}
	store := mock.NewStore()
	store.UserRepo.UpdateUserProfileFn = func(_ context.Context, _ string, _ *models.UpdateProfileRequest) error {
		return nil
	}

	body, _ := json.Marshal(models.UpdateProfileRequest{Name: "New Name"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, user))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateProfile_StorageNotFound(t *testing.T) {
	user := &models.User{Id: "user-1", Role: "STUDENT", Email: "me@example.com"}
	store := mock.NewStore()
	store.UserRepo.UpdateUserProfileFn = func(_ context.Context, _ string, _ *models.UpdateProfileRequest) error {
		return storage.ErrorNotFound
	}

	body, _ := json.Marshal(models.UpdateProfileRequest{Name: "X"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, user))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUpdatePassword_Success(t *testing.T) {
	hashedPw, _ := utils.HashPassword("old-password")
	user := &models.User{Id: "user-1", Role: "STUDENT", Email: "me@example.com"}
	store := mock.NewStore()
	store.UserRepo.UpdatePasswordFn = func(_ context.Context, _ string, _, _ string) error {
		return nil
	}
	// GetById is called inside UpdatePassword to verify the current password
	store.UserRepo.GetByIdFn = func(_ context.Context, _ *models.PrimaryKey) (*models.User, error) {
		return &models.User{Id: "user-1", Password: hashedPw}, nil
	}

	body, _ := json.Marshal(models.UpdatePasswordRequest{
		CurrentPassword: "old-password",
		NewPassword:     "new-password",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, user))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdatePassword_BadJSON(t *testing.T) {
	user := &models.User{Id: "user-1", Role: "STUDENT", Email: "me@example.com"}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile/password", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, user))

	w := newRecorder()
	newTestRouter(mock.NewStore()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
