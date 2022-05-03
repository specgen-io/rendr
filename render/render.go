package render

import (
	"fmt"
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

func renderFile(templateFile *File, argsValues blueprint.ArgsValues) (*File, error) {
	templatePath := templateFile.Path

	renderedPath, err := renderPath(templatePath, argsValues)
	if err != nil {
		return nil, err
	}

	if renderedPath == nil {
		return nil, nil
	}

	content, err := render(templateFile.Content, argsValues)
	if err != nil {
		return nil, err
	}

	return &File{*renderedPath, content, templateFile.Executable}, nil
}

func renderFiles(templateFiles []File, argsValues blueprint.ArgsValues) ([]File, error) {
	result := []File{}
	for _, templateFile := range templateFiles {
		renderedFile, err := renderFile(&templateFile, argsValues)
		if err != nil {
			return nil, fmt.Errorf(`template "%s" returned error: %s`, templateFile.Path, err.Error())
		}
		if renderedFile != nil {
			result = append(result, *renderedFile)
		}
	}
	return result, nil
}
