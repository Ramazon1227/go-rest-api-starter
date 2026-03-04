package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Ramazon1227/go-rest-api-starter/internal/model"
	"github.com/Ramazon1227/go-rest-api-starter/internal/service"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/response"
)

type AuthHandler struct {
	svc *service.UserService
}

func NewAuthHandler(svc *service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	res, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.StatusCode, svcErr.Message)
			return
		}
		response.InternalError(c)
		return
	}

	response.Created(c, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	res, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.StatusCode, svcErr.Message)
			return
		}
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}
