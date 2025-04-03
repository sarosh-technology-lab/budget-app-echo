package main

func (app *Application) routes() {

	// api routes

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

	// web routes

	webGroup := app.server.Group("/web")
	{
		webGroup.GET("/index", app.handler.ViewIndexPage)
		webGroup.GET("/admin/login", app.handler.ViewAdminLoginPage)
		webGroup.GET("/user/login", app.handler.ViewUserLoginPage)
		webGroup.GET("/category-form-page", app.handler.CategoryFormPage)
		webGroup.POST("/save-category", app.handler.SaveCategoryForm)
	}
}