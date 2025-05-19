package main

import (
	"net/http"

	"github.com/ASH-WIN-10/uniref-app-backend/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) createFileHandler(c echo.Context) error {
	clientId, err := app.readIDParam(c)
	if err != nil {
		app.notFoundResponse(c)
		return nil
	}

	file := new(data.File)
	file.Category = c.FormValue("category")
	file.ClientID = clientId

	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		app.badRequestResponse(c, err)
		return nil
	}

	err = data.SaveFileLocally(fileHeader, file)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}

	files := []data.File{*file}
	err = app.models.Files.Insert(files)
	if err != nil {
		app.serverErrorResponse(c, err)
		return nil
	}

	return c.JSONPretty(http.StatusCreated, envelope{"file": files[0]}, "\t")
}
