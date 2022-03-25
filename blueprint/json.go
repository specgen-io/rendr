package blueprint

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func ValidateValues(args Args, values map[string]interface{}) (ArgsValues, error) {
	rootArg := Map("", "", nil, args)
	value, err := ValidateValue([]string{}, &rootArg, values)
	if err != nil {
		return nil, err
	}
	return value.(ArgsValues), nil
}

func ValidateValue(path []string, arg *NamedArg, value interface{}) (interface{}, error) {
	if arg.String != nil {
		stringValue, isString := value.(string)
		if !isString {
			return nil, errors.New(fmt.Sprintf(`argument "%s" should be string`, strings.Join(path, ".")))
		}
		return stringValue, nil
	}
	if arg.Array != nil {
		arrayValues, isArray := value.([]interface{})
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
		mapValues, isMap := value.(map[string]interface{})
		if !isMap {
			return nil, errors.New(fmt.Sprintf(`argument "%s" should be map`, strings.Join(path, ".")))
		}
		values := ArgsValues{}
		for nestedArgName, nestedArgValue := range mapValues {
			nestedArg := arg.Map.Keys.FindByName(nestedArgName)
			nestedPath := append(path, nestedArg.Name)
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

func ReadValuesJson(args Args, data []byte) (ArgsValues, error) {
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
