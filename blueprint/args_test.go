package blueprint

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
	"testing"
)

var casesArgUnmarshal = []ArgUnmarshalTestCase{
	{
		"string arg",
		`
type: string
description: the description
values:
  - the value 1
  - the value 2
default: the value 1
`,
		StringArg("the description", false, "", []string{"the value 1", "the value 2"}, StrPtr("the value 1")),
	},
	{
		"bool arg",
		`
type: boolean
description: the description
default: yes
`,
		BooleanArg("the description", false, "", BoolPtr(true)),
	},
	{
		"array arg",
		`
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
		ArrayArg("the description", false, "", []string{"the value 1", "the value 2", "the value 3"}, []string{"the value 1", "the value 2"}),
	},
	{
		"group",
		`
type: group
description: the description
args:
  param:
    type: string
    description: param description
`,
		GroupArg("the description", false, "", Args{
			NamedStringArg("param", "param description", false, "", nil, nil),
		}),
	},
}

var casesArgsUnmarshal = []ArgsUnmarshalTestCase{
	{
		"three args",
		`
the_arg_1:
  type: string
  description: the description
  values:
    - the value 1
    - the value 2
  default: the value 1
the_arg_2:
  type: boolean
  description: the description
  default: yes
the_arg_3:
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
		Args{
			NamedStringArg("the_arg_1", "the description", false, "", []string{"the value 1", "the value 2"}, StrPtr("the value 1")),
			NamedBooleanArg("the_arg_2", "the description", false, "", BoolPtr(true)),
			NamedArrayArg("the_arg_3", "the description", false, "", []string{"the value 1", "the value 2", "the value 3"}, []string{"the value 1", "the value 2"}),
		},
	},
}

func Test_ArgUnmarshal(t *testing.T) {
	for _, testcase := range casesArgUnmarshal {
		t.Logf(`Running test case: %s`, testcase.Name)
		arg := Arg{}
		err := yaml.Unmarshal([]byte(strings.TrimSpace(testcase.Yaml)), &arg)
		if err != nil {
			t.Fatalf(`unmarshaling failes: %s`, err.Error())
		}
		if !cmp.Equal(testcase.Expected, arg) {
			t.Errorf("Failed, values do not match\nexpected: %v\nactual:   %v", testcase.Expected, arg)
		}
	}
}

type ArgUnmarshalTestCase struct {
	Name     string
	Yaml     string
	Expected Arg
}

func Test_ArgsUnmarshal(t *testing.T) {
	for _, testcase := range casesArgsUnmarshal {
		t.Logf(`Running test case: %s`, testcase.Name)
		args := Args{}
		err := yaml.Unmarshal([]byte(strings.TrimSpace(testcase.Yaml)), &args)
		if err != nil {
			t.Fatalf(`unmarshaling failes: %s`, err.Error())
		}
		if !cmp.Equal(testcase.Expected, args) {
			t.Errorf("Failed, values do not match\nexpected: %v\nactual:   %v", testcase.Expected, args)
		}
	}
}

type ArgsUnmarshalTestCase struct {
	Name     string
	Yaml     string
	Expected Args
}
