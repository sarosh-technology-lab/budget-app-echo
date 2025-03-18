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
	"gopkg.in/natefinch/lumberjack.v2"
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

// Configure lumberjack logger
func setupLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   "logs/app.log", // Log file path
		MaxSize:    10,             // Max size in megabytes before rotation
		MaxBackups: 3,              // Max number of old log files to keep
		MaxAge:     28,             // Max number of days to retain old log files
		Compress:   true,           // Compress rotated log files
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err.Error())
	}

	// Set up lumberjack logger
	lumberjackLogger := setupLogger()
	defer lumberjackLogger.Close()

	db, err := common.Mysql()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	e := echo.New()

	// Configure Echo logger to write to lumberjack
	logWriter := io.MultiWriter(os.Stdout, lumberjackLogger) // Logs to both console and file
	e.Logger.SetOutput(logWriter)
	// Configure middleware logger to use the same output
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Output: logWriter,
	// }))
	
	appMailer := mailer.AppMailer(e.Logger)

	// Load templates
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
		fullName := filepath.ToSlash(file)
		content, err := os.ReadFile(file)
		if err != nil {
			e.Logger.Fatalf("Failed to read %s: %v", file, err)
		}
		tmpl, err = tmpl.New(fullName).Parse(string(content))
		if err != nil {
			e.Logger.Fatalf("Failed to parse %s: %v", file, err)
		}
	}

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

	app.routes()
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	e.Logger.Print(e.Start(appAddress))
}
