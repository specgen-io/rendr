package blueprint

import (
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	"testing"
)

var casesGetValues = []GetValuesTestCase{
	{
		"string args",
		Args{
			String("param1", "", false, nil, nil),
			String("param2", "", true, nil, StrPtr("the default")),
		},
		false,
		HardcodedGetter("the value"),
		ArgsValues{"param1": "the value", "param2": "the default"},
	},
	{
		"string arg noinput",
		Args{
			String("param", "", true, nil, StrPtr("the default")),
		},
		false,
		HardcodedGetter("the value"),
		ArgsValues{"param": "the default"},
	},
	{
		"bool args",
		Args{
			Bool("param1", "", false, nil),
			Bool("param2", "", true, BoolPtr(false)),
		},
		false,
		HardcodedGetter(true),
		ArgsValues{"param1": true, "param2": false},
	},
	{
		"string args force input",
		Args{
			String("param1", "", false, nil, nil),
			String("param2", "", false, nil, StrPtr("the default")),
		},
		true,
		HardcodedGetter("the value"),
		ArgsValues{"param1": "the value", "param2": "the value"},
	},
	{
		"array args",
		Args{
			Array("param1", "", false, nil, nil),
			Array("param2", "", true, nil, []string{"three", "four"}),
		},
		false,
		HardcodedGetter([]string{"one", "two"}),
		ArgsValues{"param1": []string{"one", "two"}, "param2": []string{"three", "four"}},
	},
	{
		"array args should get",
		Args{
			Array("param1", "", false, nil, nil),
			Array("param2", "", false, nil, []string{"three", "four"}),
		},
		true,
		HardcodedGetter([]string{"one", "two"}),
		ArgsValues{"param1": []string{"one", "two"}, "param2": []string{"one", "two"}},
	},
	{
		"map args",
		Args{
			Map("themap", "", false, nil, Args{
				String("param1", "", false, nil, nil),
				String("param2", "", true, nil, StrPtr("the default")),
			}),
		},
		false,
		HardcodedGetter("the value"),
		ArgsValues{"themap": ArgsValues{"param1": "the value", "param2": "the default"}},
	},
}

func Test_GetValues(t *testing.T) {
	for _, testcase := range casesGetValues {
		t.Logf(`Running test case: %s`, testcase.Name)
		values, err := GetValues(testcase.Args, testcase.ForceInput, ArgsValues{}, testcase.Getter)
		assert.Equal(t, err, nil)
		if !cmp.Equal(testcase.Expected, values) {
			t.Errorf("\nexpected: %s\nactual:   %s", testcase.Expected, values)
		}
	}
}

type GetValuesTestCase struct {
	Name       string
	Args       Args
	ForceInput bool
	Getter     ArgValueGetter
	Expected   ArgsValues
}

func HardcodedGetter(value ArgValue) ArgValueGetter {
	return func(arg NamedArg) (ArgValue, error) {
		return value, nil
	}
}
