package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type envelope map[string]any

func readIDParam(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return 0, echo.NewHTTPError(http.StatusNotFound, "Invalid ID.")
	}
	return id, nil
}
