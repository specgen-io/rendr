package values

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/specgen-io/rendr/blueprint"
	"gotest.tools/v3/assert"
	"strings"
	"testing"
)

var casesReadValuesJson = []ReadValuesJsonTestCase{
	{
		"flat string args",
		blueprint.Args{
			blueprint.NamedStringArg("param1", "", false, "", nil, nil),
			blueprint.NamedStringArg("param2", "", false, "", nil, blueprint.StrPtr("the default")),
		},
		`{"param1":"value1","param2":"value2"}`,
		nil,
		ArgsValues{"param1": "value1", "param2": "value2"},
	},
	{
		"flat boolean args",
		blueprint.Args{
			blueprint.NamedBooleanArg("param1", "", false, "", nil),
			blueprint.NamedBooleanArg("param2", "", false, "", blueprint.BoolPtr(true)),
		},
		`{"param1":true,"param2":false}`,
		nil,
		ArgsValues{"param1": true, "param2": false},
	},
	{
		"flat string arg wrong value",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, "", nil, nil),
		},
		`{"param":123}`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"flat string arg null value",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, "", nil, nil),
		},
		`{"param":null}`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"array string arg",
		blueprint.Args{
			blueprint.NamedArrayArg("param", "", false, "", nil, nil),
		},
		`{"param":["value1","value2"]}`,
		nil,
		ArgsValues{"param": []string{"value1", "value2"}},
	},
	{
		"array string arg wrong value",
		blueprint.Args{
			blueprint.NamedArrayArg("param", "", false, "", nil, nil),
		},
		`{"param":"should be string array"}`,
		errors.New(`argument "param" should be array`),
		nil,
	},
	{
		"nested arg",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested", "", false, "", nil, nil),
			}),
		},
		`{"param":{"nested":"the_value"}}`,
		nil,
		ArgsValues{"param": ArgsValues{"nested": "the_value"}},
	},
	{
		"nested arg wrong value",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested", "", false, "", nil, nil),
			}),
		},
		`{"param":"the_value"}`,
		errors.New(`argument "param" should be map`),
		nil,
	},
	{
		"double nested arg",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedGroupArg("internal", "", false, "", blueprint.Args{
					blueprint.NamedStringArg("nested", "", false, "", nil, nil),
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
	Args     blueprint.Args
	Json     string
	Error    error
	Expected ArgsValues
}
