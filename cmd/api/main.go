package main

import (
	"budget-backend/cmd/api/handlers"
	middlewares "budget-backend/cmd/api/middleware"
	"budget-backend/common"
	"budget-backend/internal/mailer"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
)

type Application struct {
	logger echo.Logger
	server *echo.Echo
	handler handlers.Handler
	appMiddleware middlewares.AppMiddleware
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

	// e.Use(middleware.Logger(), middleware.Recover())
	app.routes()

	fmt.Print(app)
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	e.Logger.Fatal(e.Start(appAddress))
}
