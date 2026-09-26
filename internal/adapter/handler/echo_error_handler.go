package handler

import (
	"adbcon/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

func EchoErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return // no action if sent response
	}

	// default
	code := echo.ErrInternalServerError.Code
	title := "Internal Server Error"
	detail := err.Error()

	if errHttp, ok := err.(*echo.HTTPError); ok {
		// echo errors
		code = errHttp.Code
		if msg, ok := errHttp.Message.(string); ok {
			detail = msg
		}
		title = http.StatusText(code)
	} else {
		// other errors
		title = http.StatusText(code)
	}

	// make ProblemDetails
	problem := api.ProblemDetails{
		Title:  &title,
		Status: &code,
		Detail: &detail,
	}

	if err := c.JSON(code, problem); err != nil {
		c.Logger().Error(err)
	}
}
