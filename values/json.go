package values

import (
	"encoding/json"
	"fmt"
	"github.com/specgen-io/rendr/blueprint"
	"strings"
)

func ValidateValues(args blueprint.Args, values map[string]interface{}) (ArgsValues, error) {
	rootArg := blueprint.NamedGroupArg("", "", false, "", args)
	value, err := ValidateValue([]string{}, &rootArg, values)
	if err != nil {
		return nil, err
	}
	return value.(ArgsValues), nil
}

func ValidateValue(path []string, arg *blueprint.NamedArg, value interface{}) (interface{}, error) {
	if arg.String != nil {
		stringValue, isString := value.(string)
		if !isString {
			return nil, fmt.Errorf(`argument "%s" should be string`, strings.Join(path, "."))
		}
		return stringValue, nil
	}
	if arg.Boolean != nil {
		boolValue, isBool := value.(bool)
		if !isBool {
			return nil, fmt.Errorf(`argument "%s" should be boolean`, strings.Join(path, "."))
		}
		return boolValue, nil
	}
	if arg.Array != nil {
		arrayValues, isArray := value.([]interface{})
		if !isArray {
			return nil, fmt.Errorf(`argument "%s" should be array`, strings.Join(path, "."))
		}
		values := make([]string, len(arrayValues))
		for index := range arrayValues {
			values[index] = arrayValues[index].(string)
		}
		return values, nil
	}
	if arg.Map != nil {
		mapValues, isMap := value.(map[string]interface{})
		if !isMap {
			return nil, fmt.Errorf(`argument "%s" should be map`, strings.Join(path, "."))
		}
		values := ArgsValues{}
		for nestedArgName, nestedArgValue := range mapValues {
			nestedPath := append(path, nestedArgName)
			nestedArg := arg.Map.Args.FindByName(nestedArgName)
			if nestedArg == nil {
				return nil, fmt.Errorf(`argument "%s" is not defined in the blueprint but has value provided for it`, strings.Join(nestedPath, "."))
			}
			nestedValue, err := ValidateValue(nestedPath, nestedArg, nestedArgValue)
			if err != nil {
				return nil, err
			}
			values[nestedArg.Name] = nestedValue
		}
		return values, nil
	}
	panic(fmt.Sprintf(fmt.Sprintf(`unknown argument kind: "%s"`, arg.Name)))
}

func ReadValuesJson(args blueprint.Args, data []byte) (ArgsValues, error) {
	values := map[string]interface{}{}
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}
	argsValues, err := ValidateValues(args, values)
	if err != nil {
		return nil, err
	}
	return argsValues, nil
}
