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
		Arg{
			String: &ArgString{
				"the description",
				[]string{"the value 1", "the value 2"},
				StrPtr("the value 1"),
			},
		},
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
		Arg{
			Array: &ArgArray{
				"the description",
				[]string{"the value 1", "the value 2", "the value 3"},
				[]string{"the value 1", "the value 2"},
			},
		},
	},
	{
		"map arg",
		`
type: map
description: the description
keys:
  param:
    type: string
    description: param description
`,
		Arg{
			Map: &ArgMap{
				"the description",
				nil,
				Args{
					String("param", "param description", nil, nil),
				},
			},
		},
	},
}

var casesArgsUnmarshal = []ArgsUnmarshalTestCase{
	{
		"two args",
		`
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
		Args{
			{
				Name: "the_arg_1",
				Arg: Arg{
					String: &ArgString{
						"the description",
						[]string{"the value 1", "the value 2"},
						StrPtr("the value 1"),
					},
				},
			},
			{
				Name: "the_arg_2",
				Arg: Arg{
					Array: &ArgArray{
						"the description",
						[]string{"the value 1", "the value 2", "the value 3"},
						[]string{"the value 1", "the value 2"},
					},
				},
			},
		},
	},
}

func Test_ArgUnmarshal(t *testing.T) {
	for _, testcase := range casesArgUnmarshal {
		t.Logf(`Running test case: %s`, testcase.Name)
		arg := Arg{}
		yaml.Unmarshal([]byte(strings.TrimSpace(testcase.Yaml)), &arg)
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
		yaml.Unmarshal([]byte(strings.TrimSpace(testcase.Yaml)), &args)
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
