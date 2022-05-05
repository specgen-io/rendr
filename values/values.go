package values

import (
	"fmt"
	"github.com/specgen-io/rendr/blueprint"
)

type ArgValue interface{}
type ArgsValues map[string]ArgValue

type ArgValueGetter func(arg blueprint.NamedArg) (ArgValue, error)

func GetValues(args blueprint.Args, forceInput bool, noInput bool, argsValues ArgsValues, getter ArgValueGetter) (ArgsValues, error) {
	values := ArgsValues{}
	for _, arg := range args {
		condition, err := computeCondition(args, values, arg.Condition)
		if err != nil {
			return nil, err
		}
		if !condition {
			continue
		}
		value, _ := argsValues[arg.Name]
		if arg.Map != nil {
			if value == nil {
				value = ArgsValues{}
			}
			mapValue, err := GetValues(arg.Map.Keys, forceInput, noInput || arg.NoInput, value.(ArgsValues), getter)
			if err != nil {
				return nil, err
			}
			value = mapValue
		} else {
			if value == nil {
				argValue, err := getValue(arg, forceInput, noInput, getter)
				if err != nil {
					return nil, err
				}
				value = argValue
			}
		}
		values[arg.Name] = value
	}
	return values, nil
}

func computeCondition(args blueprint.Args, values ArgsValues, condition string) (bool, error) {
	result, err := RenderShort(condition, EnrichValues(args, values))
	if err != nil {
		return false, err
	}
	return result != nil, nil
}

func getValue(arg blueprint.NamedArg, forceInput bool, noInput bool, getter ArgValueGetter) (ArgValue, error) {
	isStringArgWithSingleOption := arg.String != nil && len(arg.String.Values) == 1
	shouldGet := (forceInput || (!noInput && !arg.NoInput)) && !isStringArgWithSingleOption
	value := defaultValue(arg)
	if shouldGet {
		return getter(arg)
	} else {
		if arg.NoInput && value == nil {
			return nil, fmt.Errorf(`argument "%s" doesn't have default value but marked as "noinput"'`, arg.Name)
		}
		return value, nil
	}
}

func defaultValue(arg blueprint.NamedArg) ArgValue {
	if arg.String != nil {
		if arg.String.Default != nil {
			return *arg.String.Default
		}
		return nil
	}
	if arg.Bool != nil {
		if arg.Bool.Default != nil {
			return *arg.Bool.Default
		}
		return nil
	}
	if arg.Array != nil {
		if arg.Array.Default != nil {
			return arg.Array.Default
		}
		return nil
	}
	if arg.Map != nil {
		return nil
	}
	panic(fmt.Sprintf(fmt.Sprintf(`unknown argument kind: "%s"`, arg.Name)))
}
