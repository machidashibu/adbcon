package handler_test

import (
	"adbcon/internal/adapter/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/labstack/echo/v4"
)

func TestEchoHandler(t *testing.T) {
	// testcase
	type testcase struct {
		name string
		sub  func(t *testing.T)
	}
	testcases := []testcase{
		{
			name: "GetMainPage",
			sub:  testGetMainPage,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, tc.sub)
	}
}

func testGetMainPage(t *testing.T) {
	// prepare
	h := handler.NewEchoHandler()
	require.NotNil(t, h)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/gui", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// testing
	require.ErrorIs(t, h.GetMainPage(c), echo.ErrNotFound)
}
