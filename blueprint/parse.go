package blueprint

import (
	"errors"
	"fmt"
	"strings"
)

func ParseValues(args Args, values []string) (ArgsValues, error) {
	rootArg := Map("", "", nil, args)
	result := ArgsValues{}
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		argValue := parts[1]
		path := strings.Split(parts[0], ".")

		err := setValue(&rootArg, result, path, argValue)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func setValue(arg *NamedArg, argsValues ArgsValues, path []string, argValue string) error {
	currentValues := argsValues
	for pathIndex := range path {
		argName := path[pathIndex]
		if arg.Map == nil {
			return errors.New(fmt.Sprintf(`argument "%s" should be map but found %s`, strings.Join(path[:pathIndex], "."), arg.Type()))
		}
		nextArg := arg.Map.Keys.FindByName(argName)
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
