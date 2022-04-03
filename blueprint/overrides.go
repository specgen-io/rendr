package blueprint

import (
	"errors"
	"fmt"
	"strings"
)

func OverrideValues(args Args, values, overrides ArgsValues) (ArgsValues, error) {
	rootArg := Map("", "", false, nil, args)
	value, err := OverrideValue([]string{}, &rootArg, values, overrides)
	if err != nil {
		return nil, err
	}
	return value.(ArgsValues), nil
}

func OverrideValue(path []string, arg *NamedArg, value, override ArgValue) (ArgValue, error) {
	if arg.String != nil {
		stringValue, isString := override.(string)
		if !isString {
			return nil, errors.New(fmt.Sprintf(`argument "%s" should be string`, strings.Join(path, ".")))
		}
		return stringValue, nil
	}
	if arg.Bool != nil {
		boolValue, isBool := override.(bool)
		if !isBool {
			return nil, errors.New(fmt.Sprintf(`argument "%s" should be boolean`, strings.Join(path, ".")))
		}
		return boolValue, nil
	}
	if arg.Array != nil {
		arrayValues, isArray := override.([]interface{})
		if !isArray {
			return nil, errors.New(fmt.Sprintf(`argument "%s" should be array`, strings.Join(path, ".")))
		}
		values := make([]string, len(arrayValues))
		for index := range arrayValues {
			values[index] = arrayValues[index].(string)
		}
		return values, nil
	}
	if arg.Map != nil {
		mapOverrides, err := castOrEmpty(path, override)
		if err != nil {
			return nil, err
		}

		mapValues, err := castOrEmpty(path, value)
		if err != nil {
			return nil, err
		}

		for nestedArgName, nestedOverrideValue := range mapOverrides {
			nestedArg := arg.Map.Keys.FindByName(nestedArgName)
			if nestedArg == nil {
				return nil, errors.New(``)
			}
			nestedPath := append(path, nestedArg.Name)
			nestedValue := mapValues[nestedArgName]
			newValue, err := OverrideValue(nestedPath, nestedArg, nestedValue, nestedOverrideValue)
			if err != nil {
				return nil, err
			}
			mapValues[nestedArgName] = newValue
		}
		return mapValues, nil
	}
	panic(fmt.Sprintf(fmt.Sprintf(`unknown argument kind: "%s"`, arg.Name)))
}

func castOrEmpty(path []string, value ArgValue) (ArgsValues, error) {
	if value == nil {
		return ArgsValues{}, nil
	}
	castedMapValues, isMap := value.(ArgsValues)
	if !isMap {
		return nil, errors.New(fmt.Sprintf(`argument "%s" should be map`, strings.Join(path, ".")))
	}
	if castedMapValues != nil {
		return castedMapValues, nil
	} else {
		return ArgsValues{}, nil
	}
}
