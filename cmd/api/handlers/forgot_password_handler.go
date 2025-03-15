package handlers

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
	"budget-backend/internal/mailer"
	"encoding/base64"
	"errors"
	"net/http"
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
			return common.SendNotFoundResponse(c, "Record not found, register with the provided email first", http.StatusNotFound)
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, please try again later")
	}

	token, err := appTokenService.GenerateResetPasswordToken(*user)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred, please try again later")
	}

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(user.Email))

	frontendUrl, err := url.Parse(payload.FrontendURL)
	if err != nil {
		return common.SendBadRequestResponse(c, "invalid frontend URL")
	}

	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", token.Token)
	frontendUrl.RawQuery = query.Encode()

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