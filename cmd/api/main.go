package main

import (
	"budget-backend/cmd/api/handlers"
	middlewares "budget-backend/cmd/api/middleware"
	"budget-backend/common"
	"budget-backend/internal/mailer"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	logger echo.Logger
	server *echo.Echo
	handler handlers.Handler
	appMiddleware middlewares.AppMiddleware
}

// TemplateRenderer is a custom html/template renderer for Echo
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := common.Mysql()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	e := echo.New()
	
	appMailer := mailer.AppMailer(e.Logger)

	// Load templates
	
	// Load templates with full paths as names
	tmpl := template.New("")
	adminFiles, err := filepath.Glob("internal/views/admin/*.html")
	if err != nil {
		e.Logger.Fatalf("Failed to glob admin templates: %v", err)
	}
	userFiles, err := filepath.Glob("internal/views/user/*.html")
	if err != nil {
		e.Logger.Fatalf("Failed to glob user templates: %v", err)
	}
	allFiles := append(adminFiles, userFiles...)

	for _, file := range allFiles {
		fullName := filepath.ToSlash(file) // e.g., "internal/views/admin/login.html"
		// Read the file content and define it with the full path as the name
		content, err := os.ReadFile(file)
		if err != nil {
			e.Logger.Fatalf("Failed to read %s: %v", file, err)
		}
		tmpl, err = tmpl.New(fullName).Parse(string(content))
		if err != nil {
			e.Logger.Fatalf("Failed to parse %s: %v", file, err)
		}
	}

	// Log loaded templates for debugging
	e.Logger.Infof("Loaded templates: %v", tmpl.DefinedTemplates())

	renderer := &TemplateRenderer{
		templates: tmpl,
	}
	e.Renderer = renderer

	h := handlers.Handler{
		DB: db,
		Logger: e.Logger,
		Mailer: appMailer,
	}

	appMiddleware := middlewares.AppMiddleware{
		DB: db,
		Logger: &e.Logger,
	}
	// initializing the application

	app := Application{
		logger: e.Logger,
		server: e,
		handler: h,
		appMiddleware: appMiddleware,
	}

	e.Use(middleware.Logger())
	app.routes()

	fmt.Print(app)
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	e.Logger.Fatal(e.Start(appAddress))
}
