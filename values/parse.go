package values

import (
	"errors"
	"fmt"
	"github.com/specgen-io/rendr/blueprint"
	"strconv"
	"strings"
)

func ParseValues(args blueprint.Args, values []string) (ArgsValues, error) {
	rootArg := blueprint.NamedGroupArg("", "", false, "", args)
	result := ArgsValues{}
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		argValue := parts[1]
		if strings.HasPrefix(argValue, `"`) && strings.HasSuffix(argValue, `"`) {
			argValue = argValue[1 : len(argValue)-1]
		}
		path := strings.Split(parts[0], ".")

		err := setValue(&rootArg, result, path, argValue)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func setValue(arg *blueprint.NamedArg, argsValues ArgsValues, path []string, argValue string) error {
	currentValues := argsValues
	for pathIndex := range path {
		argName := path[pathIndex]
		if arg.Map == nil {
			return errors.New(fmt.Sprintf(`argument "%s" should be map but found %s`, strings.Join(path[:pathIndex], "."), arg.Type()))
		}
		nextArg := arg.Map.Args.FindByName(argName)
		if nextArg == nil {
			return errors.New(fmt.Sprintf(`argument "%s" was not found`, strings.Join(path[:pathIndex+1], ".")))
		}
		arg = nextArg

		if pathIndex == len(path)-1 {
			if arg.Array != nil {
				argValues := strings.Split(argValue, ",")
				currentValues[argName] = argValues
			}
			if arg.String != nil {
				currentValues[argName] = argValue
			}
			if arg.Boolean != nil {
				boolValue, err := parseBoolean(argValue)
				if err != nil {
					return err
				}
				currentValues[argName] = boolValue
			}
			return nil
		} else {
			newCurrentValues, found := currentValues[argName]
			if !found {
				newCurrentValues = ArgsValues{}
				currentValues[argName] = newCurrentValues
			}
			currentValues = newCurrentValues.(ArgsValues)
		}
	}
	return nil
}

func parseBoolean(value string) (bool, error) {
	yesNoValue := parseYesNo(value)
	if yesNoValue != nil {
		return *yesNoValue, nil
	} else {
		parsedValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, err
		}
		return parsedValue, nil
	}
}

func parseYesNo(value string) *bool {
	lowerValue := strings.ToLower(value)
	if lowerValue == "yes" || lowerValue == "y" {
		return blueprint.BoolPtr(true)
	}
	if lowerValue == "no" || lowerValue == "n" {
		return blueprint.BoolPtr(false)
	}
	return nil
}
