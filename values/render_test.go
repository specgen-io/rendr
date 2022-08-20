package values

import (
	"gotest.tools/v3/assert"
	"testing"
)

var casesRenderShort = []RenderShortTestCase{
	{
		"no formula",
		"some string",
		ArgsValues{},
		StrPtr("some string"),
		nil,
	},
	{
		"no formula empty string",
		"",
		ArgsValues{},
		StrPtr(""),
		nil,
	},
	{
		"var formula",
		"{{myparam.value}}",
		ArgsValues{
			"myparam": map[string]interface{}{"value": "thevalue", "thevalue": true},
		},
		StrPtr("thevalue"),
		nil,
	},
	{
		"if formula true with string",
		"{{#myparam.thevalue}}some string",
		ArgsValues{
			"myparam": map[string]interface{}{"value": "thevalue", "thevalue": true},
		},
		StrPtr("some string"),
		nil,
	},
	{
		"if formula false with string",
		"{{#myparam.thevalue}}some string",
		ArgsValues{
			"myparam": map[string]interface{}{"value": "thevalue", "thevalue": false},
		},
		nil,
		nil,
	},
	{
		"if formula true no string",
		"{{#myparam.thevalue}}",
		ArgsValues{
			"myparam": map[string]interface{}{"value": "thevalue", "thevalue": true},
		},
		StrPtr(""),
		nil,
	},
	{
		"if formula false no string",
		"{{#myparam.thevalue}}",
		ArgsValues{
			"myparam": map[string]interface{}{"value": "thevalue", "thevalue": false},
		},
		nil,
		nil,
	},
}

func Test_RenderShort(t *testing.T) {
	for _, testcase := range casesRenderShort {
		t.Logf(`Running test case: %s`, testcase.Name)
		result, err := RenderShort(testcase.Template, testcase.Values)
		if testcase.Error != nil {
			if err == nil {
				t.Fatalf(`Expected error but got nil`)
			}
			assert.Error(t, err, testcase.Error.Error())
		} else {
			assert.NilError(t, err)
		}
		if testcase.Expected != nil {
			if result == nil {
				t.Fatalf(`Expected "%s" but got nil`, *testcase.Expected)
			}
			assert.Equal(t, *result, *testcase.Expected)
		} else {
			if result != nil {
				t.Fatalf(`Expected nil but got "%s"`, *result)
			}
		}
	}
}

type RenderShortTestCase struct {
	Name     string
	Template string
	Values   ArgsValues
	Expected *string
	Error    error
}
