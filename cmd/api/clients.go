package main

import (
	"errors"
	"fmt"
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
		CompanyName: input.CompanyName,
		ClientName:  input.ClientName,
		Email:       input.Email,
		Phone:       input.Phone,
	}

	err := app.models.Clients.Insert(&client)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/v1/clients/%d", client.ID))

	return c.JSONPretty(http.StatusCreated, envelope{"client": client}, "\t")
}

func (app *application) showClientHandler(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return nil
	}

	client, err := app.models.Clients.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return nil
	}

	return c.JSONPretty(http.StatusOK, envelope{"client": client}, "\t")
}
