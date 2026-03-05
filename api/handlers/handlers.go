package handlers

import (
	"strconv"

	"github.com/Ramazon1227/go-rest-api-starter/api/http"
	"github.com/Ramazon1227/go-rest-api-starter/config"
	"github.com/Ramazon1227/go-rest-api-starter/storage"

	"github.com/Ramazon1227/go-rest-api-starter/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg     config.Config
	log     logger.LoggerI
	storage storage.StorageI
}

func NewHandler(cfg config.Config, log logger.LoggerI, svcs storage.StorageI) Handler {
	return Handler{
		cfg:     cfg,
		log:     log,
		storage: svcs,
	}
}

func (h *Handler) handleResponse(c *gin.Context, status http.Status, data interface{}) {
	switch code := status.Code; {
	case code < 300:
		h.log.Info(
			"---Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	case code < 400:
		h.log.Warn(
			"!!!Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	default:
		h.log.Error(
			"!!!Response--->",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	}

	c.JSON(status.Code, http.Response{
		Status:      status.Status,
		Description: status.Description,
		Data:        data,
	})
}

func (h *Handler) getOffsetParam(c *gin.Context) (offset int, err error) {
	offsetStr := c.DefaultQuery("offset", h.cfg.DefaultOffset)
	return strconv.Atoi(offsetStr)
}

func (h *Handler) getLimitParam(c *gin.Context) (offset int, err error) {
	offsetStr := c.DefaultQuery("limit", h.cfg.DefaultLimit)
	return strconv.Atoi(offsetStr)
}
