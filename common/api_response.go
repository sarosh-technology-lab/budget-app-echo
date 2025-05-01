package common

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Error string `json:"error"`
}

type JSONSuccessResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type JSONVaildationErrorResponse struct {
	Status bool                       `json:"status"`
	Errors []*ValidationError `json:"errors"`
	Message string
}

func SendSuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, JSONSuccessResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SendErrorResponse(c echo.Context, message string, statusCode int) error {
	return c.JSON(statusCode, JSONErrorResponse{
		Status:  false,
		Message: message,
	})
}

func SendValidationErrorResponse(c echo.Context, errors []*ValidationError) error {
	return c.JSON(http.StatusUnprocessableEntity, JSONVaildationErrorResponse{
		Status: false,
		Message: "Vaildation Failed",
		Errors: errors,
	})
}

func SendBadRequestResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusBadRequest)
}

func SendNotFoundResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusNotFound)
}

func SendForbiddenResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusForbidden)
}

func SendInternalServerErrorResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusInternalServerError)
}
