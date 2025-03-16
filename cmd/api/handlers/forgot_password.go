package handlers

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
	"budget-backend/internal/mailer"
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) ForgotPassword(c echo.Context) error {
	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	// bind data or in simple lang retrieve the data form the request

	payload := new(requests.ForgotPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	user, err := userService.GetUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Record not found, register with the provided email first")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, please try again later")
	}

	token, err := appTokenService.GenerateResetPasswordToken(*user)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred, please try again later")
	}

	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(user.Email))

	frontendUrl, err := url.Parse(payload.FrontendURL)
	if err != nil {
		return common.SendBadRequestResponse(c, "invalid frontend URL")
	}

	query := url.Values{}
	query.Set("token", token.Token)
	query.Set("email", encodedEmail)
	frontendUrl.RawQuery = query.Encode()

	print(frontendUrl.String())

	mailData := mailer.EmailData{
		Subject: "Request Password Reset " + os.Getenv("APP_NAME"),
		Meta: struct {
			Token string
			FrontendUrl string
		}{
			Token: token.Token,
			FrontendUrl: frontendUrl.String(),
		},
	}

	// sending a welcome email to user

	err = h.Mailer.Send(payload.Email, "forgot-password.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}

	return common.SendSuccessResponse(c, "Forgot Password Email Sent", nil)
}

func (h *Handler) ResetPassword(c echo.Context) error {
	payload := new(requests.ResetPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validating the data

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	email, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		return common.SendBadRequestResponse(c, "an error occured, try again later")
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	user, err := userService.GetUserByEmail(string(email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "invalid password reset token")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, please try again later")
	}

	appToken, err := appTokenService.ValidatePasswordToken(*user, payload.Token)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	err = userService.ChangeUserPassword(*user, payload.Password)

	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	//invalidate the token, setting used = true
	appTokenService.InvalidatedToken(user.ID, *appToken)

	return common.SendSuccessResponse(c, "password is successfuly reset", nil)
}