package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/Ramazon1227/go-rest-api-starter/internal/model"
	"github.com/Ramazon1227/go-rest-api-starter/internal/repository"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
)

type UserService struct {
	repo *repository.UserRepository
	jwt  *jwt.Manager
}

func NewUserService(repo *repository.UserRepository, jwt *jwt.Manager) *UserService {
	return &UserService{repo: repo, jwt: jwt}
}

type ServiceError struct {
	StatusCode int
	Message    string
}

func (e *ServiceError) Error() string { return e.Message }

func serviceErr(status int, msg string) *ServiceError {
	return &ServiceError{StatusCode: status, Message: msg}
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, serviceErr(http.StatusConflict, "email already in use")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repo.Create(ctx, req.Email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &model.AuthResponse{Token: token, User: user}, nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, serviceErr(http.StatusUnauthorized, "invalid credentials")
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, serviceErr(http.StatusUnauthorized, "invalid credentials")
	}

	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &model.AuthResponse{Token: token, User: user}, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, serviceErr(http.StatusNotFound, "user not found")
	}
	return user, err
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req *model.UpdateProfileRequest) (*model.User, error) {
	if req.Email == "" {
		return s.repo.FindByID(ctx, userID)
	}
	user, err := s.repo.UpdateEmail(ctx, userID, req.Email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, serviceErr(http.StatusNotFound, "user not found")
	}
	return user, err
}
