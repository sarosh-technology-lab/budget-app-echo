package handlers

import (
	"budget-backend/internal/mailer"
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	Logger echo.Logger
	Mailer mailer.Mailer
}

func (h *Handler) BindBodyRequest(c echo.Context, payload any) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return errors.New("failed to bind body, make sure you are sending a valid payload")
	}
	return nil
}