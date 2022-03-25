package blueprint

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
