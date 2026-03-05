package mock

import (
	"context"

	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
)

// UserRepo implements storage.UserRepoImpl.
// Set only the function fields needed for each test; unset fields return zero values.
type UserRepo struct {
	AddFn                  func(ctx context.Context, entity *models.UserCreateModel) (*models.PrimaryKey, error)
	UpdateProfileFn        func(ctx context.Context, entity *models.UpdateUserProfileModel) error
	GetByIdFn              func(ctx context.Context, pKey *models.PrimaryKey) (*models.User, error)
	GetListFn              func(ctx context.Context, queryParam *models.QueryParam) (*models.GetUserListModel, error)
	DeleteFn               func(ctx context.Context, pKey *models.PrimaryKey) error
	GetByEmailFn           func(ctx context.Context, email string) (*models.User, error)
	UpdateUserProfileFn    func(ctx context.Context, userId string, req *models.UpdateProfileRequest) error
	UpdatePasswordFn       func(ctx context.Context, userId string, currentPassword, newPassword string) error
}

func (m *UserRepo) Add(ctx context.Context, entity *models.UserCreateModel) (*models.PrimaryKey, error) {
	if m.AddFn != nil {
		return m.AddFn(ctx, entity)
	}
	return nil, nil
}

func (m *UserRepo) UpdateProfile(ctx context.Context, entity *models.UpdateUserProfileModel) error {
	if m.UpdateProfileFn != nil {
		return m.UpdateProfileFn(ctx, entity)
	}
	return nil
}

func (m *UserRepo) GetById(ctx context.Context, pKey *models.PrimaryKey) (*models.User, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, pKey)
	}
	return nil, nil
}

func (m *UserRepo) GetList(ctx context.Context, queryParam *models.QueryParam) (*models.GetUserListModel, error) {
	if m.GetListFn != nil {
		return m.GetListFn(ctx, queryParam)
	}
	return &models.GetUserListModel{}, nil
}

func (m *UserRepo) Delete(ctx context.Context, pKey *models.PrimaryKey) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, pKey)
	}
	return nil
}

func (m *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *UserRepo) UpdateUserProfile(ctx context.Context, userId string, req *models.UpdateProfileRequest) error {
	if m.UpdateUserProfileFn != nil {
		return m.UpdateUserProfileFn(ctx, userId, req)
	}
	return nil
}

func (m *UserRepo) UpdatePassword(ctx context.Context, userId string, currentPassword, newPassword string) error {
	if m.UpdatePasswordFn != nil {
		return m.UpdatePasswordFn(ctx, userId, currentPassword, newPassword)
	}
	return nil
}

// Store implements storage.StorageI backed by a UserRepo mock.
type Store struct {
	UserRepo *UserRepo
}

func NewStore() *Store {
	return &Store{UserRepo: &UserRepo{}}
}

func (s *Store) CloseDB() {}

func (s *Store) User() storage.UserRepoImpl {
	return s.UserRepo
}
