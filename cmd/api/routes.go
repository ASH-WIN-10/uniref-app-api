package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *application) registerRoutes(e *echo.Echo) {
	e.HTTPErrorHandler = app.customHTTPErrorHandler

	e.Use(app.recoverPanic)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1420", "tauri://localhost"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Static("/assets", "assets")

	e.GET("/v1/healthcheck", app.healthcheckHandler)
	e.GET("/v1/clients", app.listClientsHandler)
	e.POST("/v1/clients", app.createClientHandler)
	e.GET("/v1/clients/:id", app.showClientHandler)
	e.DELETE("/v1/clients/:id", app.deleteClientHandler)
}
