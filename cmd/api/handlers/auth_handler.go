package handlers

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
	"budget-backend/internal/mailer"
	"budget-backend/internal/models"
	"errors"
	"os"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) Register(c echo.Context) error {
	//bind request body

	payload := new(requests.RegisterUserRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	//validation
	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	//check if email already exists

	_, err := userService.GetUserByEmail(payload.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendBadRequestResponse(c, "Email has already been taken")
	}

	registerUser, err := userService.RegisterUser(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	mailData := mailer.EmailData{
		Subject: "welcome to " + os.Getenv("APP_NAME"),
		Meta: struct{
			FirstName string
			LoginLink string
		}{
			FirstName: registerUser.FirstName,
			LoginLink: "#",
		},
	}
	// sending a welcome email to user

	err = h.Mailer.Send(payload.Email, "welcome.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}

	return common.SendSuccessResponse(c, "User Registered", registerUser)
}

func (h *Handler) Login(c echo.Context) error {
	userService := services.NewUserService(h.DB)

	// bind data or in simple lang retrieve the data form the request

	payload := new(requests.LoginRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	// check if the user with supplied email exists

	user, err := userService.GetUserByEmail(payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendBadRequestResponse(c, "Invalid credentials")
	}

	// compare the password with hashed password

	if !common.CheckPasswordHash(payload.Password, user.Password){
		return common.SendBadRequestResponse(c, "Invalid credentials")
	}

	// send response with user token

	accessToken, refreshToken, err := common.GenerateJWT((*user))
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "User Login Succesful", map[string]interface{}{
		"access_token": accessToken,
		"refresh_token": refreshToken,
		"user_id": user.ID,
	})
}

func (h *Handler) GetAuthenticationUser(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		common.SendInternalServerErrorResponse(c, "User authuentication failed")
	}
	return common.SendSuccessResponse(c, "authenticated user retreived", user)
}
