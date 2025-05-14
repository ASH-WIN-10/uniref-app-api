package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(c echo.Context, status int, message any) {
	env := envelope{"error": message}

	err := c.JSONPretty(status, env, "\t")
	if err != nil {
		app.logError(c.Request(), err)
		c.NoContent(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(c echo.Context, err error) {
	app.logError(c.Request(), err)

	message := "the server encounterd a problem and could not process your request"
	app.errorResponse(c, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(c echo.Context) {
	message := "the requested resource could not be found"
	app.errorResponse(c, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(c echo.Context) {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request().Method)
	app.errorResponse(c, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(c echo.Context, err error) {
	app.errorResponse(c, http.StatusBadRequest, err.Error())
}

// Custom HTTP error handler for Echo
func (app *application) customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		if he.Code == http.StatusMethodNotAllowed {
			app.methodNotAllowedResponse(c)
			return
		}
		if he.Code == http.StatusNotFound {
			app.notFoundResponse(c)
			return
		}
		if he.Code == http.StatusInternalServerError {
			app.serverErrorResponse(c, err)
			return
		}
	}
}
