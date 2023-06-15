package render

import (
	"github.com/specgen-io/rendr/blueprint"
	"github.com/specgen-io/rendr/input"
	"github.com/specgen-io/rendr/values"
)

func (t Template) GetArgsValues(args blueprint.Args, inputMode InputMode, valuesData *values.ValuesData, overridesKeysValues []string) (values.ArgsValues, error) {
	var err error = nil

	argsValues, err := values.ReadValuesData(args, valuesData)
	if err != nil {
		return nil, err
	}

	if argsValues == nil {
		argsValues = values.ArgsValues{}
	}

	if overridesKeysValues != nil {
		overridesValues, err := values.ParseValues(args, overridesKeysValues)
		if err != nil {
			return nil, err
		}
		argsValues, err = values.OverrideValues(args, argsValues, overridesValues)
		if err != nil {
			return nil, err
		}
	}

	argsInput := input.Survey
	if inputMode == NoInputMode {
		argsInput = input.NoInput
	}
	argsValues, err = values.GetValues(args, inputMode == ForceInputMode, inputMode == NoInputMode, argsValues, argsInput)
	if err != nil {
		return nil, err
	}

	argsValues = values.EnrichValues(args, argsValues)

	return argsValues, nil
}
