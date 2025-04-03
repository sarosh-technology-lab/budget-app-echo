package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"budget-backend/helpers"
)

func (h *Handler) ViewIndexPage(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "Welcome to My App",
		"Message": "This is the home page.",
	}

	names := map[string]interface{}{
		"Name": helpers.Capitalize("john"), // using helper functions (just for example)
		"Age":  30,
		"City": "New York",
	}

	templateData := map[string]interface{}{
		"Names": names,
		"Data": data,
	}
	
	return c.Render(http.StatusOK, "internal/views/index.html", templateData)
}
