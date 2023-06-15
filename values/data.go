package values

import (
	"encoding/json"
	"fmt"
	"github.com/specgen-io/rendr/blueprint"
	"gopkg.in/specgen-io/yaml.v3"
	"io/ioutil"
	"strings"
)

func validateValuesData(args blueprint.Args, values map[string]interface{}) (ArgsValues, error) {
	rootArg := blueprint.NamedGroupArg("", "", false, "", args)
	value, err := validateValueData([]string{}, &rootArg, values)
	if err != nil {
		return nil, err
	}
	return value.(ArgsValues), nil
}

func validateValueData(path []string, arg *blueprint.NamedArg, value interface{}) (interface{}, error) {
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
			nestedValue, err := validateValueData(nestedPath, nestedArg, nestedArgValue)
			if err != nil {
				return nil, err
			}
			values[nestedArg.Name] = nestedValue
		}
		return values, nil
	}
	panic(fmt.Sprintf(fmt.Sprintf(`unknown argument kind: "%s"`, arg.Name)))
}

type ValuesDataKind string

const (
	JSON ValuesDataKind = "json"
	YAML ValuesDataKind = "yaml"
)

type ValuesData struct {
	Kind ValuesDataKind
	Data []byte
}

func ReadValuesData(args blueprint.Args, valuesData *ValuesData) (ArgsValues, error) {
	if valuesData == nil {
		return nil, nil
	}
	values := map[string]interface{}{}
	switch valuesData.Kind {
	case JSON:
		err := json.Unmarshal(valuesData.Data, &values)
		if err != nil {
			return nil, err
		}
		break
	case YAML:
		err := yaml.Unmarshal(valuesData.Data, &values)
		if err != nil {
			return nil, err
		}
		break
	}
	argsValues, err := validateValuesData(args, values)
	if err != nil {
		return nil, err
	}
	return argsValues, nil
}

func LoadValuesFile(valuesFilePath string) (*ValuesData, error) {
	var valuesData *ValuesData = nil
	if valuesFilePath != "" {
		data, err := ioutil.ReadFile(valuesFilePath)
		if err != nil {
			return nil, fmt.Errorf(`can't open file "%s": %s`, valuesFilePath, err.Error())
		}
		valuesDataKind := JSON
		if strings.HasSuffix(valuesFilePath, ".yaml") || strings.HasSuffix(valuesFilePath, ".yml") {
			valuesDataKind = YAML
		}
		valuesData = &ValuesData{valuesDataKind, data}
	}
	return valuesData, nil
}
