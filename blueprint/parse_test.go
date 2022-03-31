package blueprint

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	"testing"
)

var casesParseValues = []ParseValuesTestCase{
	{
		"flat string args",
		Args{
			String("param1", "", nil, nil),
			String("param2", "", nil, StrPtr("the default")),
		},
		[]string{"param1=value1", "param2=value2"},
		nil,
		ArgsValues{"param1": "value1", "param2": "value2"},
	},
	{
		"flat bool args",
		Args{
			Bool("param1", "", nil),
			Bool("param2", "", BoolPtr(false)),
		},
		[]string{"param1=yes", "param2=true"},
		nil,
		ArgsValues{"param1": true, "param2": true},
	},
	{
		"non existing arg",
		Args{
			String("param1", "", nil, nil),
			String("param2", "", nil, StrPtr("the default")),
		},
		[]string{"param1=value1", "non_existing=value2"},
		errors.New(`argument "non_existing" was not found`),
		nil,
	},
	{
		"nested arg",
		Args{
			Map("param", "", nil, Args{
				String("nested", "", nil, nil),
			}),
		},
		[]string{"param.nested=the_value"},
		nil,
		ArgsValues{"param": ArgsValues{"nested": "the_value"}},
	},
	{
		"nested arg is not map",
		Args{
			String("param", "", nil, nil),
		},
		[]string{"param.nested=the_value"},
		errors.New(`argument "param" should be map but found string`),
		nil,
	},
	{
		"array arg",
		Args{
			Array("param", "", nil, nil),
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
	Args     Args
	Values   []string
	Error    error
	Expected ArgsValues
}
