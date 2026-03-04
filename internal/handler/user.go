package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/Ramazon1227/go-rest-api-starter/internal/middleware"
	"github.com/Ramazon1227/go-rest-api-starter/internal/model"
	"github.com/Ramazon1227/go-rest-api-starter/internal/service"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/response"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.StatusCode, svcErr.Message)
			return
		}
		response.InternalError(c)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.StatusCode, svcErr.Message)
			return
		}
		response.InternalError(c)
		return
	}

	response.Success(c, user)
}
