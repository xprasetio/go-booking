package routes

import (
	categoryHandler "booking/internal/category"
	facilityHandler "booking/internal/facility"
	spaceHandler "booking/internal/space"
	userHandler "booking/internal/user"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	userHandler *userHandler.UserHandler,
	categoryHandler *categoryHandler.CategoryHandler,
	spaceHandler *spaceHandler.SpaceHandler,
	facilityHandler *facilityHandler.FacilityHandler,
	authMiddleware echo.MiddlewareFunc,
	adminMiddleware echo.MiddlewareFunc,
) {
	// Public routes
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	// Protected routes
	protected := e.Group("")
	protected.Use(authMiddleware)
	{
		// User routes
		protected.POST("/logout", userHandler.Logout)
		// users routes
		users := protected.Group("/admin/v1/user")
		{
			users.GET("/me", userHandler.GetMe)
			users.PUT("/update", userHandler.UpdateProfile)
			users.DELETE("/delete", userHandler.DeleteAccount)
			users.GET("", userHandler.GetAllUsers)
			// Endpoint untuk update role (hanya superadmin)
			users.PUT("/:id/update", userHandler.UpdateUserRole)
		}
		// Category routes
		categories := protected.Group("/admin/v1/categories")
		categories.Use(adminMiddleware)
		{
			categories.POST("", categoryHandler.Create)
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}
		// Space routes
		spaces := protected.Group("/admin/v1/spaces")
		spaces.Use(adminMiddleware)
		{
			spaces.POST("", spaceHandler.Create)
			spaces.GET("", spaceHandler.GetAll)
			spaces.GET("/:id", spaceHandler.GetByID)
			spaces.PUT("/:id", spaceHandler.Update)
			spaces.DELETE("/:id", spaceHandler.Delete)
		}
		// Facility routes
		facilities := protected.Group("/admin/v1/facilities")
		facilities.Use(adminMiddleware)
		{
			facilities.POST("", facilityHandler.Create)
			facilities.GET("", facilityHandler.GetAll)
			facilities.GET("/:id", facilityHandler.GetByID)
			facilities.PUT("/:id", facilityHandler.Update)
			facilities.DELETE("/:id", facilityHandler.Delete)
		}
	}
}
