package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	httpapi "github.com/Ramazon1227/go-rest-api-starter/api/http"
	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
)

// CreateUser godoc
// @ID create-user
// @Router /v1/user [POST]
// @Summary Create User
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UserCreateModel true "user data"
// @Success 201 {object} models.User
// @Failure 400 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// @Security ApiKeyAuth
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.UserCreateModel

	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, err)
		return
	}

	resp, err := h.storage.User().Add(context.Background(), &user)
	if err != nil {
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}

	status:= httpapi.Created
	status.Description = "user created and user password has been sent to email"
	h.handleResponse(c, status, resp)
}

// GetUserByID godoc
// @ID get-user-by-id
// @Router /v1/user/{user_id} [GET]
// @Summary Get User By ID
// @Description Get user details by ID
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} httpapi.Response
// @Failure 204 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// @Security ApiKeyAuth
func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		h.handleResponse(c, httpapi.BadRequest, "user id required")
		return
	}

	resp, err := h.storage.User().GetById(context.Background(), &models.PrimaryKey{Id: id})
	if err != nil {
		if err == storage.ErrorNotFound {
			h.handleResponse(c, httpapi.NoContent, err)
			return
		}
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}

	h.handleResponse(c, httpapi.OK, resp)
}

// GetUserList godoc
// @ID get-user-list
// @Router /v1/user [GET]
// @Summary Get User List
// @Description Get list of users with pagination
// @Tags user
// @Accept json
// @Produce json
// @Param offset query integer false "offset"
// @Param limit query integer false "limit"
// @Success 200 {object} models.GetUserListModel
// @Failure 400 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// @Security ApiKeyAuth
func (h *Handler) GetUserList(c *gin.Context) {
	var queryParams models.QueryParam

	offset, err := h.getOffsetParam(c)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, "invalid offset")
		return
	}
	queryParams.Offset = int32(offset)

	limit, err := h.getLimitParam(c)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, "invalid limit")
		return
	}
	queryParams.Limit = int32(limit)

	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}

	resp, err := h.storage.User().GetList(context.Background(), &queryParams)
	if err != nil {
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}

	h.handleResponse(c, httpapi.OK, resp)
}

// UpdateUser godoc
// @ID update-user
// @Router /v1/user/{user_id} [PUT]
// @Summary Update User
// @Description Update user profile
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param user body models.UpdateUserProfileModel true "user data"
// @Success 200 {object} httpapi.Response
// @Failure 400 {object} httpapi.Response
// @Failure 204 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// @Security ApiKeyAuth
func (h *Handler) UpdateUser(c *gin.Context) {
	var user models.UpdateUserProfileModel

	user.Id = c.Param("user_id")
	if user.Id == "" {
		h.handleResponse(c, httpapi.BadRequest, "user id required")
		return
	}

	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.handleResponse(c, httpapi.BadRequest, err)
		return
	}

	err = h.storage.User().UpdateProfile(context.Background(), &user)
	if err != nil {
		if err == storage.ErrorNotFound {
			h.handleResponse(c, httpapi.NoContent, err)
			return
		}
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}

	h.handleResponse(c, httpapi.OK, "successfully updated")
}

// DeleteUser godoc
// @ID delete-user
// @Router /v1/user/{user_id} [DELETE]
// @Summary Delete User
// @Description Delete a user by ID
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} httpapi.Response
// @Failure 400 {object} httpapi.Response
// @Failure 204 {object} httpapi.Response
// @Failure 500 {object} httpapi.Response
// @Security ApiKeyAuth
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		h.handleResponse(c, httpapi.BadRequest, "user id required")
		return
	}

	err := h.storage.User().Delete(context.Background(), &models.PrimaryKey{Id: id})
	if err != nil {
		if err == storage.ErrorNotFound {
			h.handleResponse(c, httpapi.NoContent, err)
			return
		}
		h.handleResponse(c, httpapi.InternalServerError, err)
		return
	}

	h.handleResponse(c, httpapi.OK, "successfully deleted")
}
