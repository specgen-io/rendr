package render

import (
	"github.com/specgen-io/rendr/blueprint"
	"github.com/specgen-io/rendr/input"
)

func (t Template) GetArgsValues(args blueprint.Args, inputMode InputMode, valuesJsonData []byte, overridesKeysValues []string) (blueprint.ArgsValues, error) {
	var err error = nil

	argsValues := blueprint.ArgsValues{}

	if valuesJsonData != nil {
		argsValues, err = blueprint.ReadValuesJson(args, valuesJsonData)
		if err != nil {
			return nil, err
		}
	}

	if overridesKeysValues != nil {
		overridesValues, err := blueprint.ParseValues(args, overridesKeysValues)
		if err != nil {
			return nil, err
		}
		argsValues, err = blueprint.OverrideValues(args, argsValues, overridesValues)
		if err != nil {
			return nil, err
		}
	}

	argsInput := input.Survey
	if inputMode == NoInputMode {
		argsInput = input.NoInput
	}
	argsValues, err = blueprint.GetValues(args, inputMode == ForceInputMode, inputMode == NoInputMode, argsValues, argsInput)
	if err != nil {
		return nil, err
	}

	argsValues = blueprint.EnrichValues(args, argsValues)

	return argsValues, nil
}
