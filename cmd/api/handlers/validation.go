package handlers

import (
	"budget-backend/common"
	"fmt"
	"reflect"
	"strings"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ValidateRequest(c echo.Context, payload interface{}) []*common.ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var errors []*common.ValidationError
	validation := validate.Struct(payload)
	validationErrors, ok := validation.(validator.ValidationErrors)
	if ok{
		reflected := reflect.ValueOf(payload)
		for _, validationError := range validationErrors{
			field, _ := reflected.Type().FieldByName(validationError.StructField())
			key := field.Tag.Get("json")

			if key == ""{
				key = strings.ToLower(validationError.StructField())
			}

			condition := validationError.Tag()

			keyToTitleCase := strings.Replace(key, "_", " ", -1)
			param := validationError.Param()

			errMessage := ""

			switch (condition) {
			case "required":
				errMessage = keyToTitleCase + " is required"
			case "email":
				errMessage = keyToTitleCase + " must be a valid email address"	
			case "min":
					errMessage = fmt.Sprintf("%s must be of atleast %s characters", keyToTitleCase, param)
			case "max":
					errMessage = fmt.Sprintf("%s must not be greater than %s characters", keyToTitleCase, param)	
			case "len":
					errMessage = fmt.Sprintf("%s lenght must be of exact %s characters", keyToTitleCase, param)
			case "alphanum":
					errMessage = fmt.Sprintf("%s must be of %s type", keyToTitleCase, strings.ToLower(condition))
			case "alpha":
					errMessage = fmt.Sprintf("%s must be of %s type", keyToTitleCase, strings.ToLower(condition))
			case "numeric":
					errMessage = fmt.Sprintf("%s must be of %s type", keyToTitleCase, strings.ToLower(condition))
			case "oneof":
					errMessage = fmt.Sprintf("%s must be one of these (%s)", keyToTitleCase, strings.Join(strings.Split(param, " "), ", "))
			}
			currentValidationError := &common.ValidationError{
				Error: errMessage,
			}
			errors = append(errors, currentValidationError)
		}
	} 
	return errors
}