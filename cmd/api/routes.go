package main

func (app *Application) routes() {
	apiGroup := app.server.Group("/api")
	publicAuthRoutes := apiGroup.Group("/auth")
	{
		publicAuthRoutes.POST("/register", app.handler.Register)
		publicAuthRoutes.POST("/login", app.handler.Login)
		publicAuthRoutes.POST("/forgot-password", app.handler.ForgotPassword)
		publicAuthRoutes.POST("/reset-password", app.handler.ResetPassword)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthenticationMiddleware)
	{
		profileRoutes.GET("/user", app.handler.GetAuthenticationUser)
		profileRoutes.PUT("/change-password", app.handler.ChangeUserPassword)
		profileRoutes.PUT("/update-user", app.handler.UpdateUser)
	}
}