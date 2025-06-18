package main

import (
	"fmt"
	"log"

	"booking/config"
	"booking/container"
	"booking/internal/booking"
	"booking/internal/category"
	"booking/internal/facility"
	"booking/internal/space"
	spacefacility "booking/internal/space_facility"
	"booking/internal/user"
	"booking/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize container
	ctn, err := container.NewContainer()
	if err != nil {
		log.Fatal("Cannot initialize container:", err)
	}

	// Get echo instance
	e := ctn.Get(container.EchoDefName).(*echo.Echo)

	// Get handlers
	userHandler := ctn.Get(container.UserHandlerDefName).(*user.UserHandler)
	categoryHandler := ctn.Get(container.CategoryHandlerDefName).(*category.CategoryHandler)
	spaceHandler := ctn.Get(container.SpaceHandlerDefName).(*space.SpaceHandler)
	facilityHandler := ctn.Get(container.FacilityHandlerDefName).(*facility.FacilityHandler)
	spaceFacilityHandler := ctn.Get(container.SpaceFacilityHandlerDefName).(*spacefacility.SpaceFacilityHandler)
	bookingHandler := ctn.Get(container.BookingHandlerDefName).(*booking.BookingHandler)

	// Get middleware
	authMiddleware := ctn.Get(container.AuthMiddlewareDefName).(echo.MiddlewareFunc)
	adminMiddleware := ctn.Get(container.AdminAuthMiddlewareDefName).(echo.MiddlewareFunc)

	// Setup routes
	routes.SetupRoutes(e, userHandler, categoryHandler, spaceHandler, facilityHandler, spaceFacilityHandler, bookingHandler, authMiddleware, adminMiddleware)

	// Get config and start server
	cfg := ctn.Get(container.ConfigDefName).(config.Config)
	port := fmt.Sprintf(":%s", cfg.ServerPort)
	if err := e.Start(port); err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
