package server

import (
	_ "embed"
	"log/slog"
	"mime"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed assets/cert.pem
var serverCert []byte

//go:embed assets/key.pem
var serverKey []byte

//go:embed assets/html/main.html
var htmlMain []byte

//go:embed assets/css/styles.css
var cssStyles []byte

type AssetsInfo struct {
	Ext  string
	Data []byte
}

type AssetsMap map[string]AssetsInfo

var assetsTable AssetsMap = AssetsMap{
	"/gui":        AssetsInfo{Ext: ".html", Data: htmlMain},
	"/styles.css": AssetsInfo{Ext: ".css", Data: cssStyles},
}

type AssetsConfig struct {
	Assets  AssetsMap
	Skipper middleware.Skipper
}

func AssetsWithConfig(config AssetsConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultStaticConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// path := strings.TrimPrefix(path.Clean(c.Request().URL.Path), "/")
			path := c.Request().URL.Path
			info, ok := config.Assets[path]
			if !ok {
				slog.Debug("asset not found", "path", path)
				return next(c)
			}

			contentType := mime.TypeByExtension(info.Ext)
			if contentType == "" {
				contentType = http.DetectContentType(info.Data)
			}

			return c.Blob(http.StatusOK, contentType, info.Data)
		}
	}
}
