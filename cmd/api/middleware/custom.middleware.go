package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func CustomMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("custom middleware is in use")
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}