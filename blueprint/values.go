package blueprint

import "fmt"

type ArgValue interface{}
type ArgsValues map[string]ArgValue

type ArgValueGetter func(arg NamedArg) (ArgValue, error)

func GetValues(args Args, forceInput bool, noInput bool, argsValues ArgsValues, getter ArgValueGetter) (ArgsValues, error) {
	values := ArgsValues{}
	for _, arg := range args {
		value, _ := argsValues[arg.Name]
		if arg.Map != nil {
			if value == nil {
				value = ArgsValues{}
			}
			mapValue, err := GetValues(arg.Map.Keys, forceInput, noInput, value.(ArgsValues), getter)
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

func getValue(arg NamedArg, forceInput bool, noInput bool, getter ArgValueGetter) (ArgValue, error) {
	if arg.Map != nil {
		return GetValues(arg.Map.Keys, forceInput, noInput, ArgsValues{}, getter)
	} else {
		value := arg.Default()
		if (!noInput && !arg.NoInput()) || forceInput {
			return getter(arg)
		} else {
			if arg.NoInput() && value == nil {
				return nil, fmt.Errorf(`argument "%s" doesn't have default value but marked as "noinput"'`, arg.Name)
			}
			return value, nil
		}
	}
}
