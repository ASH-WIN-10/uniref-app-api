package main

import "github.com/labstack/echo/v4"

func (app *application) registerRoutes(e *echo.Echo) {
	e.GET("/v1/healthcheck", app.healthcheckHandler)

	e.POST("/v1/clients", app.createClientHandler)
	e.GET("/v1/clients/:id", app.showClientHandler)
}
