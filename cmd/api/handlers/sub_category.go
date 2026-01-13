package handlers

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
	"budget-backend/internal/custom_app_errors"
	"budget-backend/internal/models"
	"errors"
	"net/http"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ListSubCategories(c echo.Context) error {
	var subCategories []*models.SubCategory
	paginator := common.NewPaginator(subCategories, c.Request(), h.DB)
	subCategoryService := services.NewSubCategoryService(h.DB)
	pagiantedSubCategory, err := subCategoryService.List(subCategories, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "sub categories listing", pagiantedSubCategory)
}

func (h *Handler) StoreSubCategory(c echo.Context) error {
	subCategoryService := services.NewSubCategoryService(h.DB)

	// bind data or in simple lang retrieve the data form the request

	payload := new(requests.SubCategoryRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	result, err := subCategoryService.Create(*payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "sub category created successfully", result)
}

func (h *Handler) UpdateSubCategory(c echo.Context) error {
	subCategoryService := services.NewSubCategoryService(h.DB)

	// bind data or in simple lang retrieve the data form the request

	payload := new(requests.SubCategoryRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	result := subCategoryService.Update(*payload)
	if result != nil {
		return common.SendInternalServerErrorResponse(c, result.Error())
	}

	return common.SendSuccessResponse(c, "sub category updated succefully", result)
}

func (h *Handler) DeleteSubCategory(c echo.Context) error {
	subCategoryService := services.NewSubCategoryService(h.DB)
	var category_id requests.SubCategoryIDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &category_id)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	err = subCategoryService.DeleteById(category_id.ID)
	if err != nil {
		if errors.Is(err, custom_app_errors.NewNotFoundError(err.Error())){
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendBadRequestResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "sub category deleted", nil)
}

// rendering sub category saving form

func (h *Handler) SubCategoryFormPage(c echo.Context) error {
	return c.Render(http.StatusOK, "internal/views/subcategories/form.html", nil)
}

// submit html form

// func (h *Handler) SaveSubCategoryForm(c echo.Context) error {
// 	subCategoryService := services.NewSubCategoryService(h.DB)

// 	// bind data or in simple lang retrieve the data form the request

// 	payload := new(requests.SubCategoryFormRequest)
// 	if err := h.BindBodyRequest(c, payload); err != nil {
// 		return common.SendBadRequestResponse(c, err.Error())
// 	}

// 	// validate the data

// 	validationErrors := h.ValidateRequest(c, *payload)

// 	if validationErrors != nil {
// 		return common.SendValidationErrorResponse(c, validationErrors)
// 	}

// 	result, err := subCategoryService.Create(*payload)
// 	if err != nil {
// 		return common.SendInternalServerErrorResponse(c, err.Error())
// 	}

// 	return common.SendSuccessResponse(c, "sub category created succefully", result)
// }