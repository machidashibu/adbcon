package server

import (
	"mime"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// AssetsInfo is a information of asset.
type AssetsInfo struct {
	// Ext is a file extention, it is used to resolves content-type.
	Ext string
	// Data is a contents of asset.
	// It easy set by go:embed.
	Data []byte
}

// AssetsMap is a map between path and asset contents.
type AssetsMap map[string]AssetsInfo

// AssetsConfig is a configuration of Assets middleware.
type AssetsConfig struct {
	// Assets is a map for assets.
	// Default is a no action.
	Assets AssetsMap
	// Skipper is a function that decives skip condition.
	// Default is a no skip.
	Skipper middleware.Skipper
}

// AssetsWithConfig provides contents of server asset (static file) using the config.
// Relate the path in URL and byte array, and responds byte array as request.
// Not supported no cofiguration interface.
func AssetsWithConfig(config AssetsConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultStaticConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// validate
			if len(config.Assets) == 0 {
				return next(c) // skip if no config or table is emoty
			}

			// get asset
			path := c.Request().URL.Path
			info, ok := config.Assets[path]
			if !ok {
				return next(c) // asset is not found
			}
			// resolvs content-type
			contentType := mime.TypeByExtension(info.Ext)
			if contentType == "" {
				contentType = http.DetectContentType(info.Data)
			}

			// return asset
			return c.Blob(http.StatusOK, contentType, info.Data)
		}
	}
}
