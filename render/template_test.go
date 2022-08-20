package render

import (
	"github.com/specgen-io/rendr/blueprint"
	"github.com/specgen-io/rendr/values"
	"gotest.tools/v3/assert"
	"testing"
)

var renderPathTestCases = []RenderPathTestCase{
	{
		"no arguments",
		"some/path",
		values.ArgsValues{},
		blueprint.StrPtr("some/path"),
	},
	{
		"argument in path",
		"some/path_{{param.value}}_middle/item",
		values.ArgsValues{
			"param": map[string]interface{}{"value": "bla"},
		},
		blueprint.StrPtr("some/path_bla_middle/item"),
	},
	{
		"conditional empty path included",
		"some/{{#param.value}}/item",
		values.ArgsValues{
			"param": map[string]interface{}{"value": true},
		},
		blueprint.StrPtr("some/item"),
	},
	{
		"conditional empty path excluded",
		"some/{{#param.value}}/item",
		values.ArgsValues{
			"param": map[string]interface{}{"value": false},
		},
		nil,
	},
	{
		"conditional non-empty path included",
		"some/{{#param.value}}path/item",
		values.ArgsValues{
			"param": map[string]interface{}{"value": true},
		},
		blueprint.StrPtr("some/path/item"),
	},
}

func Test_RenderPath(t *testing.T) {
	for _, testcase := range renderPathTestCases {
		t.Logf(`Running test case: %s`, testcase.Name)
		renderedPath, err := renderPath(testcase.TemplatePath, testcase.ArgsValues)
		assert.Equal(t, err, nil)
		expected := "nil"
		actual := "nil"
		if testcase.Expected != nil {
			expected = *testcase.Expected
		}
		if renderedPath != nil {
			actual = *renderedPath
		}
		if actual != expected {
			t.Errorf("Failed, rendered path does not match\nexpected: %s\nactual:   %s", expected, actual)
		}
	}
}

type RenderPathTestCase struct {
	Name         string
	TemplatePath string
	ArgsValues   values.ArgsValues
	Expected     *string
}
