package blueprint

import (
	"fmt"
)

func EnrichValues(args Args, values ArgsValues) ArgsValues {
	rootArg := Map("", "", nil, args)
	value := EnrichValue(&rootArg, values)
	return value.(ArgsValues)
}

func EnrichValue(arg *NamedArg, value interface{}) interface{} {
	if arg.String != nil {
		stringValue, _ := value.(string)
		return packStringValue(arg.String.Values, stringValue)
	}
	if arg.Array != nil {
		arrayValues, _ := value.([]string)
		return packStringArrayValue(arg.Array.Values, arrayValues)
	}
	if arg.Map != nil {
		mapValues, _ := value.(ArgsValues)
		values := ArgsValues{}
		for nestedArgName, nestedArgValue := range mapValues {
			nestedArg := arg.Map.Keys.FindByName(nestedArgName)
			nestedValue := EnrichValue(nestedArg, nestedArgValue)
			values[nestedArg.Name] = nestedValue
		}
		return values
	}
	panic(fmt.Sprintf(fmt.Sprintf(`unknown argument kind: "%s"`, arg.Name)))
}

func packStringValue(possibleValues []string, stringValue string) map[string]interface{} {
	valueObj := map[string]interface{}{"value": stringValue}
	if possibleValues != nil {
		for _, possibleValue := range possibleValues {
			valueObj[possibleValue] = stringValue == possibleValue
		}
	}
	return valueObj
}

func packStringArrayValue(possibleValues []string, strArrayValue []string) map[string]interface{} {
	valueObj := map[string]interface{}{"value": strArrayValue}
	if possibleValues != nil {
		for _, possibleValue := range possibleValues {
			valueObj[possibleValue] = contains(strArrayValue, possibleValue)
		}
	}
	return valueObj
}

func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}
