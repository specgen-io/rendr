package values

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/specgen-io/rendr/blueprint"
	"gotest.tools/assert"
	"testing"
)

var casesParseValues = []ParseValuesTestCase{
	{
		"flat string args",
		blueprint.Args{
			blueprint.String("param1", "", false, nil, nil),
			blueprint.String("param2", "", false, nil, blueprint.StrPtr("the default")),
		},
		[]string{"param1=value1", "param2=value2"},
		nil,
		ArgsValues{"param1": "value1", "param2": "value2"},
	},
	{
		"flat bool args",
		blueprint.Args{
			blueprint.Bool("param1", "", false, nil),
			blueprint.Bool("param2", "", false, blueprint.BoolPtr(false)),
		},
		[]string{"param1=yes", "param2=true"},
		nil,
		ArgsValues{"param1": true, "param2": true},
	},
	{
		"non existing arg",
		blueprint.Args{
			blueprint.String("param1", "", false, nil, nil),
			blueprint.String("param2", "", false, nil, blueprint.StrPtr("the default")),
		},
		[]string{"param1=value1", "non_existing=value2"},
		errors.New(`argument "non_existing" was not found`),
		nil,
	},
	{
		"nested arg",
		blueprint.Args{
			blueprint.Map("param", "", false, blueprint.Args{
				blueprint.String("nested", "", false, nil, nil),
			}),
		},
		[]string{"param.nested=the_value"},
		nil,
		ArgsValues{"param": ArgsValues{"nested": "the_value"}},
	},
	{
		"nested arg is not map",
		blueprint.Args{
			blueprint.String("param", "", false, nil, nil),
		},
		[]string{"param.nested=the_value"},
		errors.New(`argument "param" should be map but found string`),
		nil,
	},
	{
		"array arg",
		blueprint.Args{
			blueprint.Array("param", "", false, nil, nil),
		},
		[]string{"param=value1,value2,value3"},
		nil,
		ArgsValues{"param": []string{"value1", "value2", "value3"}},
	},
}

func Test_ParseValues(t *testing.T) {
	for _, testcase := range casesParseValues {
		t.Logf(`Running test case: %s`, testcase.Name)
		values, err := ParseValues(testcase.Args, testcase.Values)
		if testcase.Error != nil {
			assert.Error(t, err, testcase.Error.Error())
		} else {
			assert.Equal(t, err, nil)
		}
		if !cmp.Equal(testcase.Expected, values) {
			t.Errorf("\nexpected: %s\nactual:   %s", testcase.Expected, values)
		}
	}
}

type ParseValuesTestCase struct {
	Name     string
	Args     blueprint.Args
	Values   []string
	Error    error
	Expected ArgsValues
}
