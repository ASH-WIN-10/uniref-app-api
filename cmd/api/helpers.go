package main

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

type envelope map[string]any

func (app *application) readIDParam(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID parameter")
	}
	return id, nil
}
