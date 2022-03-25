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
			String("param1", "", nil, nil),
			String("param2", "", nil, StrPtr("the default")),
		},
		false,
		HardcodedGetter("the value"),
		ArgsValues{"param1": "the value", "param2": "the default"},
	},
	{
		"string args should get",
		Args{
			String("param1", "", nil, nil),
			String("param2", "", nil, StrPtr("the default")),
		},
		true,
		HardcodedGetter("the value"),
		ArgsValues{"param1": "the value", "param2": "the value"},
	},
	{
		"array args",
		Args{
			Array("param1", "", nil, nil),
			Array("param2", "", nil, []string{"three", "four"}),
		},
		false,
		HardcodedGetter([]string{"one", "two"}),
		ArgsValues{"param1": []string{"one", "two"}, "param2": []string{"three", "four"}},
	},
	{
		"array args should get",
		Args{
			Array("param1", "", nil, nil),
			Array("param2", "", nil, []string{"three", "four"}),
		},
		true,
		HardcodedGetter([]string{"one", "two"}),
		ArgsValues{"param1": []string{"one", "two"}, "param2": []string{"one", "two"}},
	},
	{
		"map args",
		Args{
			Map("themap", "", nil, Args{
				String("param1", "", nil, nil),
				String("param2", "", nil, StrPtr("the default")),
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
		values, err := GetValues(testcase.Args, testcase.ReviewDefaults, ArgsValues{}, testcase.Getter)
		assert.Equal(t, err, nil)
		if !cmp.Equal(testcase.Expected, values) {
			t.Errorf("\nexpected: %s\nactual:   %s", testcase.Expected, values)
		}
	}
}

type GetValuesTestCase struct {
	Name           string
	Args           Args
	ReviewDefaults bool
	Getter         ArgValueGetter
	Expected       ArgsValues
}

func HardcodedGetter(getValue ArgValue) ArgValueGetter {
	return func(arg NamedArg) (ArgValue, error) {
		return getValue, nil
	}
}
