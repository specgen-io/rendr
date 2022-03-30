package render

import (
	"github.com/cbroglie/mustache"
	"github.com/specgen-io/rendr/blueprint"
)

func render(template string, argsValues blueprint.ArgsValues) (string, error) {
	mustache.AllowMissingVariables = false
	content, err := mustache.Render(template, argsValues)
	if err != nil {
		return "", err
	}
	return content, nil
}
