package values

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/specgen-io/rendr/blueprint"
	"gotest.tools/v3/assert"
	"strings"
	"testing"
)

var casesReadValuesJson = []ReadValuesTestCase{
	{
		"json flat string args",
		blueprint.Args{
			blueprint.NamedStringArg("param1", "", false, "", nil, nil),
			blueprint.NamedStringArg("param2", "", false, "", nil, blueprint.StrPtr("the default")),
		},
		JSON,
		`{"param1":"value1","param2":"value2"}`,
		nil,
		ArgsValues{"param1": "value1", "param2": "value2"},
	},
	{
		"json flat boolean args",
		blueprint.Args{
			blueprint.NamedBooleanArg("param1", "", false, "", nil),
			blueprint.NamedBooleanArg("param2", "", false, "", blueprint.BoolPtr(true)),
		},
		JSON,
		`{"param1":true,"param2":false}`,
		nil,
		ArgsValues{"param1": true, "param2": false},
	},
	{
		"json flat string arg wrong value",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, "", nil, nil),
		},
		JSON,
		`{"param":123}`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"json flat string arg null value",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, "", nil, nil),
		},
		JSON,
		`{"param":null}`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"json array string arg",
		blueprint.Args{
			blueprint.NamedArrayArg("param", "", false, "", nil, nil),
		},
		JSON,
		`{"param":["value1","value2"]}`,
		nil,
		ArgsValues{"param": []string{"value1", "value2"}},
	},
	{
		"json array string arg wrong value",
		blueprint.Args{
			blueprint.NamedArrayArg("param", "", false, "", nil, nil),
		},
		JSON,
		`{"param":"should be string array"}`,
		errors.New(`argument "param" should be array`),
		nil,
	},
	{
		"json nested arg",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested", "", false, "", nil, nil),
			}),
		},
		JSON,
		`{"param":{"nested":"the_value"}}`,
		nil,
		ArgsValues{"param": ArgsValues{"nested": "the_value"}},
	},
	{
		"json nested arg wrong value",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested", "", false, "", nil, nil),
			}),
		},
		JSON,
		`{"param":"the_value"}`,
		errors.New(`argument "param" should be map`),
		nil,
	},
	{
		"json double nested arg",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedGroupArg("internal", "", false, "", blueprint.Args{
					blueprint.NamedStringArg("nested", "", false, "", nil, nil),
				}),
			}),
		},
		JSON,
		`{"param":{"internal":{"nested":"the_value"}}}`,
		nil,
		ArgsValues{"param": ArgsValues{"internal": ArgsValues{"nested": "the_value"}}},
	},
}

var casesReadValuesYaml = []ReadValuesTestCase{
	{
		"yaml flat string args",
		blueprint.Args{
			blueprint.NamedStringArg("param1", "", false, "", nil, nil),
			blueprint.NamedStringArg("param2", "", false, "", nil, blueprint.StrPtr("the default")),
		},
		YAML,
		`
param1: value1
param2: value2
`,
		nil,
		ArgsValues{"param1": "value1", "param2": "value2"},
	},
	{
		"yaml flat boolean args",
		blueprint.Args{
			blueprint.NamedBooleanArg("param1", "", false, "", nil),
			blueprint.NamedBooleanArg("param2", "", false, "", blueprint.BoolPtr(true)),
		},
		YAML,
		`
param1: true
param2: false
`,
		nil,
		ArgsValues{"param1": true, "param2": false},
	},
	{
		"yaml flat string arg wrong value",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, "", nil, nil),
		},
		YAML,
		`
param: 123
`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"yaml flat string arg null value",
		blueprint.Args{
			blueprint.NamedStringArg("param", "", false, "", nil, nil),
		},
		YAML,
		`
param: null
`,
		errors.New(`argument "param" should be string`),
		nil,
	},
	{
		"yaml array string arg",
		blueprint.Args{
			blueprint.NamedArrayArg("param", "", false, "", nil, nil),
		},
		YAML,
		`
param: [value1, value2]
`,
		nil,
		ArgsValues{"param": []string{"value1", "value2"}},
	},
	{
		"yaml array string arg wrong value",
		blueprint.Args{
			blueprint.NamedArrayArg("param", "", false, "", nil, nil),
		},
		YAML,
		`
param: should be string array
`,
		errors.New(`argument "param" should be array`),
		nil,
	},
	{
		"yaml nested arg",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested", "", false, "", nil, nil),
			}),
		},
		YAML,
		`
param:
  nested: the_value
`,
		nil,
		ArgsValues{"param": ArgsValues{"nested": "the_value"}},
	},
	{
		"yaml nested arg wrong value",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedStringArg("nested", "", false, "", nil, nil),
			}),
		},
		YAML,
		`
param: the_value
`,
		errors.New(`argument "param" should be map`),
		nil,
	},
	{
		"yaml double nested arg",
		blueprint.Args{
			blueprint.NamedGroupArg("param", "", false, "", blueprint.Args{
				blueprint.NamedGroupArg("internal", "", false, "", blueprint.Args{
					blueprint.NamedStringArg("nested", "", false, "", nil, nil),
				}),
			}),
		},
		YAML,
		`
param:
  internal: {nested: the_value }
`,
		nil,
		ArgsValues{"param": ArgsValues{"internal": ArgsValues{"nested": "the_value"}}},
	},
}

func ExecuteReadValuesTestCases(t *testing.T, testCases []ReadValuesTestCase) {
	for _, testcase := range testCases {
		t.Logf(`Running test case: %s`, testcase.Name)
		values, err := ReadValuesData(testcase.Args, &ValuesData{testcase.Kind, []byte(strings.TrimSpace(testcase.Data))})
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

func Test_ReadValues(t *testing.T) {
	ExecuteReadValuesTestCases(t, casesReadValuesJson)
	ExecuteReadValuesTestCases(t, casesReadValuesYaml)
}

type ReadValuesTestCase struct {
	Name     string
	Args     blueprint.Args
	Kind     ValuesDataKind
	Data     string
	Error    error
	Expected ArgsValues
}
