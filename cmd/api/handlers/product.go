package handlers

import (
	"budget-backend/cmd/api/requests"
	"budget-backend/cmd/api/services"
	"budget-backend/common"
	"budget-backend/internal/custom_app_errors"
	"budget-backend/internal/models"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "net/http"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ListProducts(c echo.Context) error {
	var products []*models.Product
	paginator := common.NewPaginator(products, c.Request(), h.DB)
	productService := services.NewProductService(h.DB)
	pagiantedProduct, err := productService.List(products, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "products listing", pagiantedProduct)
}

func (h *Handler) StoreProduct(c echo.Context) error {
	productService := services.NewProductService(h.DB)

	// Handle image file first
	file, err := c.FormFile("image")
	var imagePath string
	if file != nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		filename := fmt.Sprintf("%d%s", time.Now().UnixMilli(), ext)
		savePath := filepath.Join("public", "images", "products", filename)

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot open file"})
		}
		defer src.Close()

		dst, err := os.Create(savePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot save file"})
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "copy failed"})
		}

		imagePath = "/images/products/" + filename
	}

	// Try binding the form data directly
	payload := new(requests.ProductRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, "Failed to bind form data: " + err.Error())
	}

	// Override image path if file was uploaded
	if imagePath != "" {
		payload.Image = imagePath
	}

	// Validate the data
	validationErrors := h.ValidateRequest(c, *payload)
	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	result, err := productService.Create(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "product created successfully", result)
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	productService := services.NewProductService(h.DB)

	// bind data or in simple lang retrieve the data form the request

	payload := new(requests.ProductRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data

	validationErrors := h.ValidateRequest(c, *payload)

	if validationErrors != nil {
		return common.SendValidationErrorResponse(c, validationErrors)
	}

	result := productService.Update(*payload)
	if result != nil {
		return common.SendInternalServerErrorResponse(c, result.Error())
	}

	return common.SendSuccessResponse(c, "product updated succefully", result)
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	productService := services.NewProductService(h.DB)
	var product_id requests.ProductIDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &product_id)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	err = productService.DeleteById(product_id.ID)
	if err != nil {
		if errors.Is(err, custom_app_errors.NewNotFoundError(err.Error())){
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendBadRequestResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "product deleted", nil)
}

// rendering sub category saving form

// func (h *Handler) ProductFormPage(c echo.Context) error {
// 	return c.Render(http.StatusOK, "internal/views/subcategories/form.html", nil)
// }

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
