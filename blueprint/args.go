package blueprint

import (
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
)

type Args []NamedArg

func (args Args) FindByName(name string) *NamedArg {
	for index := range args {
		if args[index].Name == name {
			return &args[index]
		}
	}
	return nil
}

func (value *Args) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := Args{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		model, err := unmarshalArg(keyNode, valueNode)
		if err != nil {
			return err
		}
		array = append(array, *model)
	}
	*value = array
	return nil
}

func unmarshalArg(keyNode *yaml.Node, valueNode *yaml.Node) (*NamedArg, error) {
	var name string
	err := keyNode.DecodeWith(decodeStrict, &name)
	if err != nil {
		return nil, err
	}
	model := Arg{}
	err = valueNode.DecodeWith(decodeStrict, &model)
	if err != nil {
		return nil, err
	}
	return &NamedArg{Name: name, Arg: model}, nil
}

func (args Args) Get(name string) *NamedArg {
	for i, arg := range args {
		if arg.Name == name {
			return &args[i]
		}
	}
	return nil
}

type NamedArg struct {
	Name string
	Arg
}

type Arg struct {
	Bool   *ArgBool
	String *ArgString
	Array  *ArgArray
	Map    *ArgMap
}

func (value *Arg) UnmarshalYAML(node *yaml.Node) error {
	arg := Arg{}
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be mapping")
	}

	typ, err := decodeStringOptional(node, "type")

	if err != nil {
		return err
	}

	if typ == nil {
		yamlError(node, `field "type" is required but missing`)
	}

	switch *typ {
	case `string`:
		argString := ArgString{}
		err := node.DecodeWith(decodeLooze, &argString)
		if err != nil {
			return err
		}
		arg.String = &argString
		break
	case `boolean`:
		argBool := ArgBool{}
		err := node.DecodeWith(decodeLooze, &argBool)
		if err != nil {
			return err
		}
		arg.Bool = &argBool
		break
	case `array`:
		argArray := ArgArray{}
		err := node.DecodeWith(decodeLooze, &argArray)
		if err != nil {
			return err
		}
		arg.Array = &argArray
		break
	case `map`:
		argMap := ArgMap{}
		err := node.DecodeWith(decodeLooze, &argMap)
		if err != nil {
			return err
		}
		arg.Map = &argMap
		break
	default:
		return yamlError(node, fmt.Sprintf(`unknown argument type: %s`, *typ))
	}

	*value = arg
	return nil
}

func (arg *NamedArg) Type() ArgType {
	if arg.String != nil {
		return ArgTypeString
	}
	if arg.Bool != nil {
		return ArgTypeBoolean
	}
	if arg.Array != nil {
		return ArgTypeArray
	}
	if arg.Map != nil {
		return ArgTypeMap
	}
	panic(fmt.Sprintf(fmt.Sprintf(`unknown argument kind: "%s"`, arg.Name)))
}

type ArgType string

const (
	ArgTypeString  ArgType = "string"
	ArgTypeBoolean ArgType = "boolean"
	ArgTypeArray   ArgType = "array"
	ArgTypeMap     ArgType = "map"
)

func (arg NamedArg) NoInput() bool {
	if arg.String != nil {
		return arg.String.NoInput
	}
	if arg.Bool != nil {
		return arg.Bool.NoInput
	}
	if arg.Array != nil {
		return arg.Array.NoInput
	}
	if arg.Map != nil {
		return false
	}
	panic(fmt.Sprintf(fmt.Sprintf(`unknown argument kind: "%s"`, arg.Name)))
}

type ArgString struct {
	Description string   `yaml:"description"`
	NoInput     bool     `yaml:"noinput"`
	Values      []string `yaml:"values"`
	Default     *string  `yaml:"default"`
}

type ArgArray struct {
	Description string   `yaml:"description"`
	NoInput     bool     `yaml:"noinput"`
	Values      []string `yaml:"values"`
	Default     []string `yaml:"default"`
}

type ArgBool struct {
	Description string `yaml:"description"`
	NoInput     bool   `yaml:"noinput"`
	Default     *bool  `yaml:"default"`
}

type ArgMap struct {
	Description string `yaml:"description"`
	NoInput     bool   `yaml:"noinput"`
	Keys        Args   `yaml:"keys"`
}

func String(name string, description string, noinput bool, values []string, defaultValue *string) NamedArg {
	return NamedArg{
		Name: name,
		Arg: Arg{
			String: &ArgString{description, noinput, values, defaultValue},
		},
	}
}

func Bool(name string, description string, noinput bool, defaultValue *bool) NamedArg {
	return NamedArg{
		Name: name,
		Arg: Arg{
			Bool: &ArgBool{description, noinput, defaultValue},
		},
	}
}

func Array(name string, description string, noinput bool, values []string, defaultValue []string) NamedArg {
	return NamedArg{
		Name: name,
		Arg: Arg{
			Array: &ArgArray{description, noinput, values, defaultValue},
		},
	}
}

func Map(name string, description string, noinput bool, keys Args) NamedArg {
	return NamedArg{
		Name: name,
		Arg: Arg{
			Map: &ArgMap{description, noinput, keys},
		},
	}
}
