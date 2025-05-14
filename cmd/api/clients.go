package main

import (
	"net/http"

	"github.com/ASH-WIN-10/Himwan-Refrigerations-Backend/internal/data"
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

	data := data.Client{
		ID:          id,
		CompanyName: "John Doe Inc.",
		ClientName:  "John Doe",
		Email:       "johndoe@example.com",
		Phone:       "123-456-7890",
	}

	return c.JSONPretty(http.StatusOK, data, "\t")
}
