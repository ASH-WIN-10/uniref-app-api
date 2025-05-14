package main

import (
	"net/http"

	"github.com/ASH-WIN-10/Himwan-Refrigerations-Backend/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) createClientHandler(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	input := new(data.Client)
	if err := c.Bind(input); err != nil {
		app.badRequestResponse(c, err)
		return nil
	}

	client := data.Client{
		ID:          input.ID,
		CompanyName: input.CompanyName,
		ClientName:  input.ClientName,
		Email:       input.Email,
		Phone:       input.Phone,
	}

	return c.JSONPretty(http.StatusOK, envelope{"client": client}, "\t")
}

func (app *application) showClientHandler(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return nil
	}

	data := data.Client{
		ID:          id,
		CompanyName: "John Doe Inc.",
		ClientName:  "John Doe",
		Email:       "johndoe@example.com",
		Phone:       "123-456-7890",
	}

	return c.JSONPretty(http.StatusOK, envelope{"client": data}, "\t")
}
