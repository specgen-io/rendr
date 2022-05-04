package values

import (
	"github.com/google/go-cmp/cmp"
	"github.com/specgen-io/rendr/blueprint"
	"gotest.tools/assert"
	"testing"
)

var casesGetValues = []GetValuesTestCase{
	{
		"string args",
		blueprint.Args{
			blueprint.String("param1", "", false, nil, nil),
			blueprint.String("param2", "", true, nil, blueprint.StrPtr("the default")),
		},
		false,
		false,
		HardcodedGetter("the value"),
		ArgsValues{"param1": "the value", "param2": "the default"},
	},
	{
		"string arg noinput",
		blueprint.Args{
			blueprint.String("param", "", true, nil, blueprint.StrPtr("the default")),
		},
		false,
		false,
		HardcodedGetter("the value"),
		ArgsValues{"param": "the default"},
	},
	{
		"bool args",
		blueprint.Args{
			blueprint.Bool("param1", "", false, nil),
			blueprint.Bool("param2", "", true, blueprint.BoolPtr(false)),
		},
		false,
		false,
		HardcodedGetter(true),
		ArgsValues{"param1": true, "param2": false},
	},
	{
		"string args force input",
		blueprint.Args{
			blueprint.String("param1", "", false, nil, nil),
			blueprint.String("param2", "", false, nil, blueprint.StrPtr("the default")),
		},
		true,
		false,
		HardcodedGetter("the value"),
		ArgsValues{"param1": "the value", "param2": "the value"},
	},
	{
		"array args",
		blueprint.Args{
			blueprint.Array("param1", "", false, nil, nil),
			blueprint.Array("param2", "", true, nil, []string{"three", "four"}),
		},
		false,
		false,
		HardcodedGetter([]string{"one", "two"}),
		ArgsValues{"param1": []string{"one", "two"}, "param2": []string{"three", "four"}},
	},
	{
		"array args should get",
		blueprint.Args{
			blueprint.Array("param1", "", false, nil, nil),
			blueprint.Array("param2", "", false, nil, []string{"three", "four"}),
		},
		true,
		false,
		HardcodedGetter([]string{"one", "two"}),
		ArgsValues{"param1": []string{"one", "two"}, "param2": []string{"one", "two"}},
	},
	{
		"map args",
		blueprint.Args{
			blueprint.Map("themap", "", false, blueprint.Args{
				blueprint.String("param1", "", false, nil, nil),
				blueprint.String("param2", "", true, nil, blueprint.StrPtr("the default")),
			}),
		},
		false,
		false,
		HardcodedGetter("the value"),
		ArgsValues{"themap": ArgsValues{"param1": "the value", "param2": "the default"}},
	},
}

func Test_GetValues(t *testing.T) {
	for _, testcase := range casesGetValues {
		t.Logf(`Running test case: %s`, testcase.Name)
		values, err := GetValues(testcase.Args, testcase.ForceInput, testcase.NoInput, ArgsValues{}, testcase.Getter)
		assert.Equal(t, err, nil)
		if !cmp.Equal(testcase.Expected, values) {
			t.Errorf("\nexpected: %s\nactual:   %s", testcase.Expected, values)
		}
	}
}

type GetValuesTestCase struct {
	Name       string
	Args       blueprint.Args
	ForceInput bool
	NoInput    bool
	Getter     ArgValueGetter
	Expected   ArgsValues
}

func HardcodedGetter(value ArgValue) ArgValueGetter {
	return func(arg blueprint.NamedArg) (ArgValue, error) {
		return value, nil
	}
}
