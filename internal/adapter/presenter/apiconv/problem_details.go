package apiconv

import (
	"adbcon/api"
	"net/http"
)

// MakeProblemDetails makes ProblemDetails pbject.
func MakeProblemDetails(code int, err error) api.ProblemDetails {
	status := code
	title := http.StatusText(status)
	detail := err.Error()

	return api.ProblemDetails{
		Status: &status,
		Title:  &title,
		Detail: &detail,
	}
}

// MakeBadRequest makes ProblemDetails object with status code = 400.
func MakeBadRequest(err error) api.ProblemDetails {
	return MakeProblemDetails(http.StatusBadRequest, err)
}

// MakeInternalServerError makes ProblemDetails object with status code = 500.
func MakeInternalServerError(err error) api.ProblemDetails {
	return MakeProblemDetails(http.StatusInternalServerError, err)
}
