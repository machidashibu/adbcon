package controller

import (
	"adbcon/api"
	"log/slog"
	"net/url"
)

func GetServerAddr() (string, error) {
	urlServer, err := url.Parse(api.ServerUrlADBConsoleServer)
	if err != nil {
		slog.Error("URL parse error", "err", err, "url", api.ServerUrlADBConsoleServer)
		return "", err
	}
	return urlServer.Host, nil
}
