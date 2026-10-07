package controller

import (
	"adbcon/api"
	"adbcon/internal/adapter/model"
	"adbcon/internal/domain"
	"strings"
)

type AdbVersionParser struct {
	app string
}

func NewAdbVersionParser(app string) *AdbVersionParser {
	return &AdbVersionParser{
		app: app,
	}
}

func (a AdbVersionParser) Parse(result domain.CommandOutput) (domain.Version, error) {
	version := model.Version{App: a.app}

	lines := result.Lines()

	// get adb version
	if len(lines) > 0 && strings.HasPrefix(lines[0], "Android Debug Bridge version") {
		version.Adb = lines[0][len("Android Debug Bridge version")+1:]
	}

	// get android sdk version
	if len(lines) > 1 && strings.HasPrefix(lines[1], "Version") {
		version.Sdk = lines[1][len("Version")+1:]
	}

	// get api version
	if spec, err := api.GetSpec(); err == nil {
		version.Api = spec.Info.Version
	}

	return version, nil
}
