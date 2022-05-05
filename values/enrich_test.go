package values

import (
	"github.com/google/go-cmp/cmp"
	"github.com/specgen-io/rendr/blueprint"
	"testing"
)

var enrichValues = []EnrichTestCase{
	{
		"string arg",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, nil, nil),
		},
		ArgsValues{"param": "the value"},
		ArgsValues{
			"param": map[string]interface{}{"value": "the value"},
		},
	},
	{
		"boolean arg",
		blueprint.Args{
			blueprint.NamedBooleanArg("param", "", false, nil),
		},
		ArgsValues{"param": true},
		ArgsValues{
			"param": map[string]interface{}{"value": true},
		},
	},
	{
		"string with values arg",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, []string{"value1", "value2"}, nil),
		},
		ArgsValues{"param": "value2"},
		ArgsValues{
			"param": map[string]interface{}{"value": "value2", "value1": false, "value2": true},
		},
	},
	{
		"array string arg",
		blueprint.Args{
			blueprint.NamedArrayArg("param", "", false, []string{"value1", "value2", "value3"}, nil),
		},
		ArgsValues{"param": []string{"value1", "value3"}},
		ArgsValues{
			"param": map[string]interface{}{"value": []string{"value1", "value3"}, "value1": true, "value2": false, "value3": true},
		},
	},
}

func Test_EnrichValues(t *testing.T) {
	for _, testcase := range enrichValues {
		t.Logf(`Running test case: %s`, testcase.Name)
		values := EnrichValues(testcase.Args, testcase.Input)
		if !cmp.Equal(testcase.Expected, values) {
			t.Errorf("Failed, values do not match\nexpected: %s\nactual:   %s", testcase.Expected, values)
		}
	}
}

type EnrichTestCase struct {
	Name     string
	Args     blueprint.Args
	Input    ArgsValues
	Expected ArgsValues
}
