package values

import (
	"github.com/google/go-cmp/cmp"
	"github.com/specgen-io/rendr/blueprint"
	"gotest.tools/v3/assert"
	"testing"
)

var casesOverrideValues = []OverrideValuesTestCase{
	{
		"string arg",
		blueprint.Args{
			blueprint.NamedStringArg("param1", "", false, "", nil, nil),
			blueprint.NamedStringArg("param2", "", false, "", nil, blueprint.StrPtr("the default")),
		},
		ArgsValues{"param1": "the value", "param2": "the default"},
		ArgsValues{"param2": "the override"},
		ArgsValues{"param1": "the value", "param2": "the override"},
	},
	{
		"boolean arg",
		blueprint.Args{
			blueprint.NamedBooleanArg("param1", "", false, "", nil),
			blueprint.NamedBooleanArg("param2", "", false, "", blueprint.BoolPtr(false)),
		},
		ArgsValues{"param1": false, "param2": false},
		ArgsValues{"param2": true},
		ArgsValues{"param1": false, "param2": true},
	},
	{
		"string arg new",
		blueprint.Args{
			blueprint.NamedStringArg("param1", "", false, "", nil, nil),
			blueprint.NamedStringArg("param2", "", false, "", nil, blueprint.StrPtr("the default")),
		},
		ArgsValues{"param1": "the value"},
		ArgsValues{"param2": "the override"},
		ArgsValues{"param1": "the value", "param2": "the override"},
	},
	{
		"nested arg",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested1", "", false, "", nil, nil),
				blueprint.NamedStringArg("nested2", "", false, "", nil, nil),
			}),
		},
		ArgsValues{"param": ArgsValues{"nested1": "the_value"}},
		ArgsValues{"param": ArgsValues{"nested2": "override"}},
		ArgsValues{"param": ArgsValues{"nested1": "the_value", "nested2": "override"}},
	},
	{
		"nested arg from nil",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested", "", false, "", nil, nil),
			}),
		},
		nil,
		ArgsValues{"param": ArgsValues{"nested": "override"}},
		ArgsValues{"param": ArgsValues{"nested": "override"}},
	},
}

func Test_OverrideValues(t *testing.T) {
	for _, testcase := range casesOverrideValues {
		t.Logf(`Running test case: %s`, testcase.Name)
		values, err := OverrideValues(testcase.Args, testcase.Values, testcase.Overrides)
		assert.Equal(t, err, nil)
		if !cmp.Equal(testcase.Expected, values) {
			t.Errorf("\nexpected: %s\nactual:   %s", testcase.Expected, values)
		}
	}
}

type OverrideValuesTestCase struct {
	Name      string
	Args      blueprint.Args
	Values    ArgsValues
	Overrides ArgsValues
	Expected  ArgsValues
}
