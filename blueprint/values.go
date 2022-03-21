package blueprint

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ArgValue interface{}
type ArgsValues map[string]ArgValue

type ArgValueGetter func(arg NamedArg) (ArgValue, error)

func GetValues(args Args, reviewDefaults bool, argsValues ArgsValues, getter ArgValueGetter) (ArgsValues, error) {
	values := ArgsValues{}
	for _, arg := range args {
		value, _ := argsValues[arg.Name]
		if value == nil {
			argValue, err := getValue(arg, reviewDefaults, getter)
			if err != nil {
				return nil, err
			}
			value = argValue
		}
		values[arg.Name] = value
	}
	return values, nil
}

func getValue(arg NamedArg, reviewDefaults bool, getter ArgValueGetter) (ArgValue, error) {
	value := arg.Default()
	if value == nil || reviewDefaults {
		if arg.Map != nil {
			if value == nil {
				value = ArgsValues{}
			}
			return GetValues(arg.Map.Keys, reviewDefaults, value.(ArgsValues), getter)
		} else {
			return getter(arg)
		}
	}
	return value, nil
}

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
		nextArg := findArgByName(arg.Map.Keys, argName)
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
			nestedArg := findArgByName(arg.Map.Keys, nestedArgName)
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

func findArgByName(args Args, name string) *NamedArg {
	for index := range args {
		if args[index].Name == name {
			return &args[index]
		}
	}
	return nil
}

func OverrideValues(args Args, values, overrides ArgsValues) (ArgsValues, error) {
	rootArg := Map("", "", nil, args)
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
		mapOverrides, isMap := override.(ArgsValues)
		if !isMap {
			return nil, errors.New(fmt.Sprintf(`argument "%s" should be map`, strings.Join(path, ".")))
		}
		mapValues, isMap := value.(ArgsValues)
		if !isMap {
			return nil, errors.New(fmt.Sprintf(`argument "%s" should be map`, strings.Join(path, ".")))
		}

		for nestedArgName, nestedOverrideValue := range mapOverrides {
			nestedArg := findArgByName(arg.Map.Keys, nestedArgName)
			if nestedArg == nil {
				return nil, errors.New(``)
			}
			nestedPath := append(path, nestedArg.Name)
			nestedValue, found := mapValues[nestedArgName]
			if !found {
				nestedValue = ArgsValues{}
			}
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
