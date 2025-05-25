	package handlers

	import (
		"budget-backend/cmd/api/requests"
		"budget-backend/cmd/api/services"
		"budget-backend/common"
		// "budget-backend/internal/mailer"
		"budget-backend/internal/models"
		"errors"
		"net/http"
		// "os"
		"time"

		"github.com/labstack/echo/v4"
		"gorm.io/gorm"
	)

	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	func (h *Handler) Register(c echo.Context) error {
		//bind request body

		payload := new(requests.RegisterUserRequest)
		if err := h.BindBodyRequest(c, payload); err != nil {
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

		if payload.Phone != "" {
			_, err := userService.GetUserByPhone(payload.Phone)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
					return common.SendBadRequestResponse(c, "Phone has already been taken")
			}
		}

		registerUser, err := userService.RegisterUser(payload)
		if err != nil {
			return common.SendInternalServerErrorResponse(c, err.Error())
		}

		// mailData := mailer.EmailData{
		// 	Subject: "welcome to " + os.Getenv("APP_NAME"),
		// 	Meta: struct {
		// 		FirstName string
		// 		LoginLink string
		// 	}{
		// 		FirstName: registerUser.FirstName,
		// 		LoginLink: "#",
		// 	},
		// }

		// // sending a welcome email to user

		// err = h.Mailer.Send(payload.Email, "welcome.html", mailData)
		// if err != nil {
		// 	h.Logger.Error(err)
		// }

		return common.SendSuccessResponse(c, "User Registered", registerUser)
	}

	func (h *Handler) Login(c echo.Context) error {
		userService := services.NewUserService(h.DB)

		// bind data or in simple lang retrieve the data form the request

		payload := new(requests.LoginRequest)
		if err := h.BindBodyRequest(c, payload); err != nil {
			return common.SendBadRequestResponse(c, err.Error())
		}

		// validate the data

		validationErrors := h.ValidateRequest(c, *payload)

		if validationErrors != nil {
			return common.SendValidationErrorResponse(c, validationErrors)
		}

		// check if the user with supplied email exists

		user, err := userService.GetUserByEmail(payload.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.SendBadRequestResponse(c, "Invalid credentials")
			}
			return common.SendInternalServerErrorResponse(c, "Something went wrong please try again later")
		}

		// compare the password with hashed password

		if !common.CheckPasswordHash(payload.Password, user.Password) {
			return common.SendBadRequestResponse(c, "Invalid credentials")
		}

		// send response with user token

		accessToken, refreshToken, err := common.GenerateJWT(*user, h.DB)
		if err != nil {
			return common.SendInternalServerErrorResponse(c, err.Error())
		}

		return common.SendSuccessResponse(c, "User Login Succesful", map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user_id":       user.ID,
			"role_id":       user.RoleId,
		})
	}

	func (h *Handler) ViewAdminLoginPage(c echo.Context) error {
		// Render the login.html template
		return c.Render(http.StatusOK, "internal/views/admin/login.html", nil)
	}

	func (h *Handler) ViewUserLoginPage(c echo.Context) error {
		// Render the login.html template
		return c.Render(http.StatusOK, "internal/views/user/login.html", nil)
	}

	func (h *Handler) GetAuthenticationUser(c echo.Context) error {
		user, ok := c.Get("user").(models.User)
		if !ok {
			return common.SendInternalServerErrorResponse(c, "User authuentication failed")
		}
		return common.SendSuccessResponse(c, "authenticated user retreived", user)
	}

	func (h *Handler) RefreshToken(c echo.Context) error {
		// Bind and validate request
		payload := new(RefreshRequest)
		if err := h.BindBodyRequest(c, payload); err != nil {
			return common.SendBadRequestResponse(c, err.Error())
		}
	
		validationErrors := h.ValidateRequest(c, *payload)
		if validationErrors != nil {
			return common.SendValidationErrorResponse(c, validationErrors)
		}
	
		// Parse and validate refresh token
		claims, err := common.ParseJWTSignedAccessToken(payload.RefreshToken)
		if err != nil {
			return common.SendBadRequestResponse(c, "Invalid refresh token")
		}
	
		// Check if token is expired
		if common.IsClaimExpired(claims) {
			return common.SendBadRequestResponse(c, "Refresh token expired")
		}
	
		// Verify refresh token exists in the database
		var refreshToken models.RefreshToken
		if err := h.DB.Where("token = ? AND user_id = ?", payload.RefreshToken, claims.ID).First(&refreshToken).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.SendBadRequestResponse(c, "Invalid refresh token")
			}
			return common.SendInternalServerErrorResponse(c, "Something went wrong, please try again later")
		}
	
		// Check if refresh token is expired in the database
		if refreshToken.ExpiresAt.Before(time.Now()) {
			return common.SendBadRequestResponse(c, "Refresh token expired")
		}
	
		// Fetch user to generate new tokens
		userService := services.NewUserService(h.DB)
		user, err := userService.GetUserByID(uint(claims.ID))
		if err != nil {
			return common.SendInternalServerErrorResponse(c, "Something went wrong, please try again later")
		}
	
		// Generate new access token (and optionally new refresh token)
		newAccessToken, newRefreshToken, err := common.GenerateJWT(*user, h.DB)
		if err != nil {
			return common.SendInternalServerErrorResponse(c, err.Error())
		}
	
		// Optionally, delete the old refresh token to enforce rotation
		if err := h.DB.Delete(&refreshToken).Error; err != nil {
			return common.SendInternalServerErrorResponse(c, "Something went wrong, please try again later")
		}
	
		// Send response
		return common.SendSuccessResponse(c, "Token refreshed successfully", map[string]interface{}{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"user_id":       user.ID,
			"role_id":       user.RoleId,
		})
	}
