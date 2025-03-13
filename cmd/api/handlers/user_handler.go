package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type user struct {
	Name string `json:name`
}

func (h *Handler) UserInfo(c echo.Context) error {
	
	userInfo := user{
		Name: "sarosh",
	}

	return c.JSON(http.StatusOK, userInfo)
}