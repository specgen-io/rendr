package values

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"strings"
)

func Render(template string, argsValues ArgsValues) (string, error) {
	mustache.AllowMissingVariables = false
	content, err := mustache.Render(template, argsValues)
	if err != nil {
		return "", err
	}
	return content, nil
}

func RenderShort(template string, argsValues ArgsValues) (*string, error) {
	if strings.HasPrefix(template, "{{#") {
		closeIndex := strings.Index(template, "}}")
		// if closeIndex == -1
		formula := template[:closeIndex+2]
		internal := template[closeIndex+2:]
		if internal == "" {
			internal = "PLACEHOLDER"
		}
		fullTemplate := closeFormula(formula, internal)

		_, err := Render(fmt.Sprintf(`{{%s}}`, getArgument(formula)), argsValues)
		if err != nil {
			return nil, err
		}

		result, err := Render(fullTemplate, argsValues)
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
		result, err := Render(template, argsValues)
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
