package main

func (app *Application) routes() {

	app.server.GET("admin-login-page", app.handler.ViewAdminLoginPage)
	app.server.GET("user-login-page", app.handler.ViewUserLoginPage)
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

	categoryRoutes := apiGroup.Group("/category", app.appMiddleware.AuthenticationMiddleware)
	{
		categoryRoutes.GET("/list", app.handler.ListCategories)
		categoryRoutes.POST("/store", app.handler.StoreCategory)
		categoryRoutes.DELETE("/delete/:id", app.handler.DeleteCategory)
	}
}