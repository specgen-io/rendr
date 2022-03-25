package blueprint

import (
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	"testing"
)

var casesOverrideValues = []OverrideValuesTestCase{
	{
		"string arg",
		Args{
			String("param1", "", nil, nil),
			String("param2", "", nil, StrPtr("the default")),
		},
		ArgsValues{"param1": "the value", "param2": "the default"},
		ArgsValues{"param2": "the override"},
		ArgsValues{"param1": "the value", "param2": "the override"},
	},
	{
		"string arg new",
		Args{
			String("param1", "", nil, nil),
			String("param2", "", nil, StrPtr("the default")),
		},
		ArgsValues{"param1": "the value"},
		ArgsValues{"param2": "the override"},
		ArgsValues{"param1": "the value", "param2": "the override"},
	},
	{
		"nested arg",
		Args{
			Map("param", "", nil, Args{
				String("nested1", "", nil, nil),
				String("nested2", "", nil, nil),
			}),
		},
		ArgsValues{"param": ArgsValues{"nested1": "the_value"}},
		ArgsValues{"param": ArgsValues{"nested2": "override"}},
		ArgsValues{"param": ArgsValues{"nested1": "the_value", "nested2": "override"}},
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
	Args      Args
	Values    ArgsValues
	Overrides ArgsValues
	Expected  ArgsValues
}
