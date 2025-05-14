package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) createClientHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Client created\n")
}

func (app *application) showClientHandler(c echo.Context) error {
	id, err := readIDParam(c)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("Show the details of client with ID: %d\n", id))
}
