package render

import (
	"fmt"
	"github.com/specgen-io/rendr/blueprint"
	"strings"
)

func renderPath(templatePath string, argsValues blueprint.ArgsValues) (*string, error) {
	parts := strings.Split(templatePath, "/")
	resultParts := []string{}
	for _, part := range parts {
		resultPart, err := renderShortTemplate(part, argsValues)
		if err != nil {
			return nil, err
		}
		if resultPart == nil {
			return nil, nil
		}
		if *resultPart != "" {
			resultParts = append(resultParts, *resultPart)
		}
	}
	result := strings.Join(resultParts, "/")
	return &result, nil
}

func renderShortTemplate(template string, argsValues blueprint.ArgsValues) (*string, error) {
	if strings.HasPrefix(template, "{{#") {
		closeIndex := strings.Index(template, "}}")
		// if closeIndex == -1
		formula := template[:closeIndex+2]
		internal := template[closeIndex+2:]
		if internal == "" {
			internal = "PLACEHOLDER"
		}
		fullTemplate := closeFormula(formula, internal)

		_, err := render(fmt.Sprintf(`{{%s}}`, getArgument(formula)), argsValues)
		if err != nil {
			return nil, err
		}

		result, err := render(fullTemplate, argsValues)
		if err != nil {
			return nil, err
		}

		if result == "" {
			return nil, nil
		}
		if result == "PLACEHOLDER" {
			result = ""
		}
		return &result, nil
	} else {
		result, err := render(template, argsValues)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
}

func closeFormula(formula string, internal string) string {
	argument := getArgument(formula)
	formulaClosing := fmt.Sprintf(`{{/%s}}`, argument)
	return fmt.Sprintf(`%s%s%s`, formula, internal, formulaClosing)
}

func getArgument(formula string) string {
	s := formula
	s = strings.TrimLeft(s, "{#^")
	s = strings.TrimRight(s, "}")
	s = strings.TrimSpace(s)
	return s
}
