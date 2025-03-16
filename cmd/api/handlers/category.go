package handlers

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
	"budget-backend/internal/custom_app_errors"
	"budget-backend/internal/models"
	"errors"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ListCategories(c echo.Context) error {
	var categories []*models.Category
	paginator := common.NewPaginator(categories, c.Request(), h.DB)
	categoryService := services.NewCategoryService(h.DB)
	pagiantedCategory, err := categoryService.List(categories, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "categories listing", pagiantedCategory)
}

func (h *Handler) StoreCategory(c echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)

	// bind data or in simple lang retrieve the data form the request

	payload := new(requests.CategoryRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	result, err := categoryService.Create(*payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "category created succefully", result)
}

func (h *Handler) DeleteCategory(c echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)
	var category_id requests.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &category_id)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	err = categoryService.DeleteById(category_id.ID)
	if err != nil {
		if errors.Is(err, custom_app_errors.NewNotFoundError(err.Error())){
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendBadRequestResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "category deleted", nil)
}