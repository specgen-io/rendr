package input

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/specgen-io/rendr/blueprint"
)

func Survey(arg blueprint.NamedArg) (blueprint.ArgValue, error) {
	if arg.String != nil {
		return getString(arg)
	}
	if arg.Bool != nil {
		return getBool(arg)
	}
	if arg.Array != nil {
		return getArray(arg)
	}
	return nil, errors.New(fmt.Sprintf(`unknown kind of argument "%s"`, arg.Name))
}

func getBool(arg blueprint.NamedArg) (blueprint.ArgValue, error) {
	defaultValue := true
	if arg.Bool.Default != nil {
		defaultValue = *arg.Bool.Default
	}
	message := fmt.Sprintf(`%s:`, arg.Bool.Description)
	value := false
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &value)
	return value, err
}

func getString(arg blueprint.NamedArg) (blueprint.ArgValue, error) {
	defaultValue := ""
	if arg.String.Default != nil {
		defaultValue = *arg.String.Default
	}
	message := fmt.Sprintf(`%s:`, arg.String.Description)
	value := ""
	var prompt survey.Prompt = nil
	if arg.String.Values != nil {
		prompt = &survey.Select{
			Message: message,
			Options: arg.String.Values,
			Default: defaultValue,
		}
	} else {
		prompt = &survey.Input{
			Message: message,
			Default: defaultValue,
		}
	}
	err := survey.AskOne(prompt, &value)
	return value, err
}

func getArray(arg blueprint.NamedArg) (blueprint.ArgValue, error) {
	message := fmt.Sprintf(`%s:`, arg.Array.Description)
	value := []string{}
	prompt := &survey.MultiSelect{
		Message: message,
		Options: arg.Array.Values,
		Default: arg.Array.Default,
	}
	err := survey.AskOne(prompt, &value)
	return value, err
}
