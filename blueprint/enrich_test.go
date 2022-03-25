package blueprint

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

var enrichValues = []EnrichTestCase{
	{
		"string arg",
		Args{
			String("param", "", nil, nil),
		},
		ArgsValues{"param": "the value"},
		ArgsValues{
			"param": map[string]interface{}{"value": "the value"},
		},
	},
	{
		"string with values arg",
		Args{
			String("param", "", []string{"value1", "value2"}, nil),
		},
		ArgsValues{"param": "value2"},
		ArgsValues{
			"param": map[string]interface{}{"value": "value2", "value1": false, "value2": true},
		},
	},
	{
		"array string arg",
		Args{
			Array("param", "", []string{"value1", "value2", "value3"}, nil),
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
	Args     Args
	Input    ArgsValues
	Expected ArgsValues
}
