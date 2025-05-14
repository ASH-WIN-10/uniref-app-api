package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (app *application) createClientHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Client created\n")
}

func (app *application) showClientHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Client not found\n")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Show the details of client with ID: %d\n", id))
}
