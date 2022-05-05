package blueprint

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
	"testing"
)

var casesBlueprintUnmarshal = []BlueprintUnmarshalTestCase{
	{
		"blueprint with args",
		`
blueprint: 0
name: sample blueprint
title: The Sample Blueprint
args:
  the_arg_1:
    type: string
    description: the description
    values:
      - the value 1
      - the value 2
    default: the value 1
  the_arg_2:
    type: array
    description: the description
    values:
      - the value 1
      - the value 2
      - the value 3
    default:
      - the value 1
      - the value 2
`,
		Blueprint{
			Blueprint: "0",
			Name:      "sample blueprint",
			Title:     "The Sample Blueprint",
			Args: Args{
				NamedStringArg("the_arg_1", "the description", false, "", []string{"the value 1", "the value 2"}, StrPtr("the value 1")),
				NamedArrayArg("the_arg_2", "the description", false, "", []string{"the value 1", "the value 2", "the value 3"}, []string{"the value 1", "the value 2"}),
			},
		},
	},
}

func Test_BlueprintUnmarshal(t *testing.T) {
	for _, testcase := range casesBlueprintUnmarshal {
		t.Logf(`Running test case: %s`, testcase.Name)
		blueprint := Blueprint{}
		yaml.Unmarshal([]byte(strings.TrimSpace(testcase.Yaml)), &blueprint)
		if !cmp.Equal(testcase.Expected, blueprint) {
			t.Errorf("Failed, values do not match\nexpected: %v\nactual:   %v", testcase.Expected, blueprint)
		}
	}
}

type BlueprintUnmarshalTestCase struct {
	Name     string
	Yaml     string
	Expected Blueprint
}
