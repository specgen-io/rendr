package blueprint

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	"strings"
	"testing"
)

var casesReadValuesJson = []ReadValuesJsonTestCase{
	{
		"flat string args",
		Args{
			String("param1", "", nil, nil),
			String("param2", "", nil, StrPtr("the default")),
		},
		`{"param1":"value1","param2":"value2"}`,
		nil,
		ArgsValues{"param1": "value1", "param2": "value2"},
	},
	{
		"flat boolean args",
		Args{
			Bool("param1", "", nil),
			Bool("param2", "", BoolPtr(true)),
		},
		`{"param1":true,"param2":false}`,
		nil,
		ArgsValues{"param1": true, "param2": false},
	},
	{
		"flat string arg wrong value",
		Args{
			String("param", "", nil, nil),
		},
		`{"param":123}`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"flat string arg null value",
		Args{
			String("param", "", nil, nil),
		},
		`{"param":null}`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"array string arg",
		Args{
			Array("param", "", nil, nil),
		},
		`{"param":["value1","value2"]}`,
		nil,
		ArgsValues{"param": []string{"value1", "value2"}},
	},
	{
		"array string arg wrong value",
		Args{
			Array("param", "", nil, nil),
		},
		`{"param":"should be string array"}`,
		errors.New(`argument "param" should be array`),
		nil,
	},
	{
		"nested arg",
		Args{
			Map("param", "", nil, Args{
				String("nested", "", nil, nil),
			}),
		},
		`{"param":{"nested":"the_value"}}`,
		nil,
		ArgsValues{"param": ArgsValues{"nested": "the_value"}},
	},
	{
		"nested arg wrong value",
		Args{
			Map("param", "", nil, Args{
				String("nested", "", nil, nil),
			}),
		},
		`{"param":"the_value"}`,
		errors.New(`argument "param" should be map`),
		nil,
	},
	{
		"double nested arg",
		Args{
			Map("param", "", nil, Args{
				Map("internal", "", nil, Args{
					String("nested", "", nil, nil),
				}),
			}),
		},
		`{"param":{"internal":{"nested":"the_value"}}}`,
		nil,
		ArgsValues{"param": ArgsValues{"internal": ArgsValues{"nested": "the_value"}}},
	},
}

func Test_ReadValues(t *testing.T) {
	for _, testcase := range casesReadValuesJson {
		t.Logf(`Running test case: %s`, testcase.Name)
		values, err := ReadValuesJson(testcase.Args, []byte(strings.TrimSpace(testcase.Json)))
		if testcase.Error != nil {
			assert.Error(t, err, testcase.Error.Error())
		} else {
			assert.Equal(t, err, nil)
		}
		if !cmp.Equal(testcase.Expected, values) {
			t.Errorf("Failed, values do not match\nexpected: %s\nactual:   %s", testcase.Expected, values)
		}
	}
}

type ReadValuesJsonTestCase struct {
	Name     string
	Args     Args
	Json     string
	Error    error
	Expected ArgsValues
}
