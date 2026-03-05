package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
	"github.com/Ramazon1227/go-rest-api-starter/tests/mock"
)

// adminUser is a helper that returns a SYSTEM_ADMIN user for protected requests.
func adminUser() *models.User {
	return &models.User{Id: "admin-1", Role: "SYSTEM_ADMIN", Email: "admin@example.com"}
}

func TestCreateUser_Success(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.AddFn = func(_ context.Context, entity *models.UserCreateModel) (*models.PrimaryKey, error) {
		return &models.PrimaryKey{Id: "new-user-id"}, nil
	}

	body, _ := json.Marshal(models.UserCreateModel{
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  "STUDENT",
		Phone: "+1234567890",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateUser_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(mock.NewStore()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_StorageError(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.AddFn = func(_ context.Context, _ *models.UserCreateModel) (*models.PrimaryKey, error) {
		return nil, fmt.Errorf("db error")
	}

	body, _ := json.Marshal(models.UserCreateModel{Name: "Bob", Email: "bob@example.com", Role: "STUDENT"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateUser_Forbidden(t *testing.T) {
	// Non-admin role should get 403
	student := &models.User{Id: "s-1", Role: "STUDENT", Email: "s@example.com"}

	body, _ := json.Marshal(models.UserCreateModel{Name: "X", Email: "x@example.com", Role: "STUDENT"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, student))

	w := newRecorder()
	newTestRouter(mock.NewStore()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUserByID_Success(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.GetByIdFn = func(_ context.Context, pKey *models.PrimaryKey) (*models.User, error) {
		return &models.User{Id: pKey.Id, Name: "Alice", Email: "alice@example.com"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/user-42", nil)
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "user-42", data["id"])
}

func TestGetUserByID_NotFound(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.GetByIdFn = func(_ context.Context, _ *models.PrimaryKey) (*models.User, error) {
		return nil, storage.ErrorNotFound
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/nonexistent", nil)
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetUserList_Success(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.GetListFn = func(_ context.Context, _ *models.QueryParam) (*models.GetUserListModel, error) {
		return &models.GetUserListModel{
			Count: 2,
			Users: []*models.User{
				{Id: "u1", Name: "Alice"},
				{Id: "u2", Name: "Bob"},
			},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user?limit=10&offset=0", nil)
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["count"])
}

func TestUpdateUser_Success(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.UpdateProfileFn = func(_ context.Context, _ *models.UpdateUserProfileModel) error {
		return nil
	}

	body, _ := json.Marshal(models.UpdateUserProfileModel{Name: "Updated"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/user-1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser_NotFound(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.UpdateProfileFn = func(_ context.Context, _ *models.UpdateUserProfileModel) error {
		return storage.ErrorNotFound
	}

	body, _ := json.Marshal(models.UpdateUserProfileModel{Name: "X"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/ghost", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUser_Success(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.DeleteFn = func(_ context.Context, _ *models.PrimaryKey) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/user/user-1", nil)
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	store := mock.NewStore()
	store.UserRepo.DeleteFn = func(_ context.Context, _ *models.PrimaryKey) error {
		return storage.ErrorNotFound
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/user/ghost", nil)
	req.Header.Set("Authorization", bearerToken(t, adminUser()))

	w := newRecorder()
	newTestRouter(store).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
