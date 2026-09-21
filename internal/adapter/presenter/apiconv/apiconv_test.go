package apiconv_test

import (
	"adbcon/internal/adapter/presenter/apiconv"
	"errors"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestProblemDetails(t *testing.T) {
	var errorTestDummy = errors.New("error test dummy")

	{ // testing: MakeProblemDetails
		api := apiconv.MakeProblemDetails(http.StatusInternalServerError, errorTestDummy)
		require.Equal(t, 500, *api.Status)
		require.Equal(t, "Internal Server Error", *api.Title)
		require.Equal(t, "error test dummy", *api.Detail)
	}
	{ // testing: MakeBadRequest
		api := apiconv.MakeBadRequest(errorTestDummy)
		require.Equal(t, 400, *api.Status)
		require.Equal(t, "Bad Request", *api.Title)
		require.Equal(t, "error test dummy", *api.Detail)
	}
	{ // testing: MakeInternalServerError
		api := apiconv.MakeInternalServerError(errorTestDummy)
		require.Equal(t, 500, *api.Status)
		require.Equal(t, "Internal Server Error", *api.Title)
		require.Equal(t, "error test dummy", *api.Detail)
	}
}
