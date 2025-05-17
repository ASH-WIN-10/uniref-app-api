package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ASH-WIN-10/uniref-app-backend/internal/data"
	"github.com/ASH-WIN-10/uniref-app-backend/internal/validator"
	"github.com/labstack/echo/v4"
)

func (app *application) createClientHandler(c echo.Context) error {
	input := new(data.Client)

	err := echo.FormFieldBinder(c).
		String("client_name", &input.ClientName).
		String("company_name", &input.CompanyName).
		String("email", &input.Email).
		String("phone", &input.Phone).
		BindError()

	if err != nil {
		app.badRequestResponse(c, err)
		return nil
	}

	form, err := c.MultipartForm()
	if err != nil {
		app.badRequestResponse(c, err)
		return nil
	}

	client := data.Client{
		CompanyName: input.CompanyName,
		ClientName:  input.ClientName,
		Email:       input.Email,
		Phone:       input.Phone,
	}

	v := validator.New()
	if data.ValidateClient(v, &client); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return nil
	}

	err = app.models.Clients.Insert(&client)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}

	filesMetadata, err := app.SaveFilesLocally(form, client.ID)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}

	err = app.models.Files.Insert(filesMetadata)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}
	client.Files = filesMetadata

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

	files, err := app.models.Files.Get(id)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}

	client.Files = files

	return c.JSONPretty(http.StatusOK, envelope{"client": client}, "\t")
}

func (app *application) deleteClientHandler(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return nil
	}

	err = app.models.Clients.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return nil
	}

	return c.JSONPretty(http.StatusOK, envelope{"message": "client successfully deleted"}, "\t")
}

func (app *application) listClientsHandler(c echo.Context) error {
	var input struct {
		CompanyName string
		data.Filters
	}

	v := validator.New()
	qs := c.QueryParams()

	input.CompanyName = app.readString(qs, "company_name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "company_name", "-id", "-company_name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return nil
	}

	clients, metadata, err := app.models.Clients.GetAll(input.CompanyName, input.Filters)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}

	return c.JSONPretty(http.StatusOK, envelope{"clients": clients, "metadata": metadata}, "\t")
}

func (app *application) updateClientHandler(c echo.Context) error {
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

	var input struct {
		CompanyName string `json:"company_name"`
		ClientName  string `json:"client_name"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
	}

	err = echo.FormFieldBinder(c).
		String("client_name", &input.ClientName).
		String("company_name", &input.CompanyName).
		String("email", &input.Email).
		String("phone", &input.Phone).
		BindError()

	if err != nil {
		app.badRequestResponse(c, err)
		return nil
	}

	client.CompanyName = input.CompanyName
	client.ClientName = input.ClientName
	client.Email = input.Email
	client.Phone = input.Phone

	v := validator.New()
	if data.ValidateClient(v, client); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return nil
	}

	err = app.models.Clients.Update(client)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return nil
	}

	return c.JSONPretty(http.StatusOK, envelope{"client": client}, "\t")
}
