package middlewares

import (
	"budget-backend/common"
	"budget-backend/internal/models"
	"errors"
	"net/http"
	"strings"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AppMiddleware struct{
	Logger *echo.Logger
	DB *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return common.SendErrorResponse(c, "Please provide bearer token", http.StatusUnauthorized)
		}
		authHeaderSplit := strings.Split(authHeader, " ")
		accessToken := authHeaderSplit[1]
		claims, err := common.ParseJWTSignedAccessToken(accessToken)
		if err != nil {
			return common.SendErrorResponse(c, err.Error(), http.StatusUnauthorized)
		}
		if common.IsClaimExpired((claims)){
			return common.SendErrorResponse(c, "Token is expired", http.StatusUnauthorized)
		}

		var user models.User
		result := appMiddleware.DB.First(&user, claims.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.SendErrorResponse(c, "Invalid access Token", http.StatusUnauthorized)
		}
		if result.Error != nil {
			return common.SendErrorResponse(c, "Invalid access Token", http.StatusUnauthorized)
		}
		c.Set("user", user)
		return next(c)
	}
}