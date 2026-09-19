package server_test

import (
	"adbcon/internal/infra/server"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/labstack/echo/v4"
)

func TestAssets(t *testing.T) {
	// test case
	const (
		contentTypeCss = "text/css"
		contentTypeJs  = "text/javascript" // application/javascript (echo.MIMEApplicationJavaScript) is obsolete by RFC
	)
	var (
		testdataHtml = []byte("<html><hader><title>/main.html</title></header></html>")
		testdataCss  = []byte("bodt { color: red; }")
		testdataJs   = []byte("const el = document.getElementByTag(`body');")

		errorSkip     = errors.New("skip")
		errorNoAssets = errors.New("assets")

		nextDoNotCall = func(c echo.Context) error { return fmt.Errorf("called next") }
		nextSkip      = func(c echo.Context) error { return errorSkip }
		nextNoAssets  = func(c echo.Context) error { return errorNoAssets }
		nextNotFound  = func(c echo.Context) error { return echo.ErrNotFound }
	)
	type testcase struct {
		name            string
		config          server.AssetsConfig
		next            echo.HandlerFunc
		reqMethod       string
		reqPath         string
		reqBody         io.Reader
		respError       error
		respCode        int
		respContentType string
		respBody        []byte
	}
	testcases := []testcase{
		{
			name: "get html",
			config: server.AssetsConfig{
				Assets: server.AssetsMap{
					"/main.html": server.AssetsInfo{
						Ext:  ".html",
						Data: testdataHtml,
					},
				},
			},
			next:            nextDoNotCall,
			reqMethod:       http.MethodGet,
			reqPath:         "/main.html",
			reqBody:         strings.NewReader(""),
			respCode:        http.StatusOK,
			respContentType: echo.MIMETextHTML,
			respBody:        testdataHtml,
		},
		{
			name: "get css",
			config: server.AssetsConfig{
				Assets: server.AssetsMap{
					"/styles.css": server.AssetsInfo{
						Ext:  ".css",
						Data: testdataCss,
					},
				},
			},
			next:            nextDoNotCall,
			reqMethod:       http.MethodGet,
			reqPath:         "/styles.css",
			reqBody:         strings.NewReader(""),
			respCode:        http.StatusOK,
			respContentType: contentTypeCss,
			respBody:        testdataCss,
		},
		{
			name: "get JavaScript",
			config: server.AssetsConfig{
				Assets: server.AssetsMap{
					"/script.js": server.AssetsInfo{
						Ext:  ".js",
						Data: testdataJs,
					},
				},
			},
			next:            nextDoNotCall,
			reqMethod:       http.MethodGet,
			reqPath:         "/script.js",
			reqBody:         strings.NewReader(""),
			respCode:        http.StatusOK,
			respContentType: contentTypeJs,
			respBody:        testdataJs,
		},
		{
			name: "skipper",
			config: server.AssetsConfig{
				Skipper: func(c echo.Context) bool { return true },
				Assets: server.AssetsMap{
					"/main.html": server.AssetsInfo{
						Ext:  ".html",
						Data: testdataHtml,
					},
				},
			},
			next:      nextSkip,
			reqMethod: http.MethodGet,
			reqPath:   "/main.html",
			reqBody:   strings.NewReader(""),
			respError: errorSkip,
		},
		{
			name:      "empty config (no assets)",
			config:    server.AssetsConfig{},
			next:      nextNoAssets,
			reqMethod: http.MethodGet,
			reqPath:   "/main.html",
			reqBody:   strings.NewReader(""),
			respError: errorNoAssets,
		},
		{
			name: "not found assets",
			config: server.AssetsConfig{
				Assets: server.AssetsMap{
					"/main.html": server.AssetsInfo{
						Ext:  ".html",
						Data: testdataHtml,
					},
				},
			},
			next:      nextNotFound,
			reqMethod: http.MethodGet,
			reqPath:   "/other.html",
			reqBody:   strings.NewReader(""),
			respError: echo.ErrNotFound,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mw := server.AssetsWithConfig(tc.config)
			require.NotNil(t, mw)

			e := echo.New()
			req := httptest.NewRequest(tc.reqMethod, tc.reqPath, tc.reqBody)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			f := mw(tc.next)
			require.NotNil(t, f)
			if tc.respError != nil {
				// check to occurred error
				require.ErrorIs(t, f(c), tc.respError)
			} else {
				// check to not occurred error and respons
				require.NoError(t, f(c))
				require.Equal(t, tc.respCode, rec.Code)
				require.Truef(t, strings.HasPrefix(rec.Header().Get(echo.HeaderContentType), tc.respContentType),
					"Not matched Content-Type, Actual: %s, Expect: %s", rec.Header().Get(echo.HeaderContentType), tc.respContentType)
				require.Equal(t, tc.respBody, rec.Body.Bytes())
			}
		})
	}
}
