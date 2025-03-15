package handlers

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
	"budget-backend/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) ChangeUserPassword(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authuentication failed")
	}

	payload := new(requests.ChangePasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	if !common.CheckPasswordHash(payload.CurrentPassword, user.Password) {
		return common.SendBadRequestResponse(c, "the supplied current password does not match your current password")
	}

	userService := services.NewUserService(h.DB)
	err := userService.ChangeUserPassword(user, payload.Password)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "password changed successfuly", nil)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authuentication failed")
	}

	payload := new(requests.UpdateUserRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	updateMap := make(map[string]any)

	if payload.FirstName != "" {
		updateMap["first_name"] = payload.FirstName
	}
	if payload.LastName != "" {
		updateMap["last_name"] = payload.LastName
	}
	if payload.Email != "" {
		// Check if email already exists, ignoring the current user
		var existingUser models.User
		err := h.DB.Where("email = ? AND id != ?", payload.Email, user.ID).First(&existingUser).Error
		if err == nil { // Email exists and belongs to another user
			return common.SendBadRequestResponse(c, "Email is already taken by another user")
		} else if err != gorm.ErrRecordNotFound { // Unexpected error
			return common.SendInternalServerErrorResponse(c, "Failed to check email availability")
		}

		updateMap["email"] = payload.Email
	}
	if payload.Gender != "" {
		updateMap["gender"] = payload.Gender
	}
	if payload.Address != "" {
		updateMap["address"] = payload.Address
	}
	if payload.Phone != "" {
		var existingUser models.User
		err := h.DB.Where("phone = ? AND id != ?", payload.Phone, user.ID).First(&existingUser).Error
		if err == nil { // Phone exists and belongs to another user
			return common.SendBadRequestResponse(c, "Phone is already taken by another user")
		} else if err != gorm.ErrRecordNotFound { // Unexpected error
			return common.SendInternalServerErrorResponse(c, "Failed to check phone availability")
		}

		updateMap["phone"] = payload.Phone
	}

	userService := services.NewUserService(h.DB)
	err := userService.UpdateUser(user, updateMap)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "user updated suucessfuly", nil)
}
