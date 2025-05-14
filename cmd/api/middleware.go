package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (app *application) recoverPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.Response().Header().Set("Connection", "close")
				app.serverErrorResponse(c, fmt.Errorf("%s", err))
			}
		}()

		return next(c)
	}
}
