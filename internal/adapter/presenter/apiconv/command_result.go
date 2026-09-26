package apiconv

import (
	"adbcon/api"
	"adbcon/internal/domain"
)

func CommandResult(result domain.CommandResult) api.CommandResult {
	return api.CommandResult{
		Serial: result.Serial(),
		Result: result.Text(),
	}
}
