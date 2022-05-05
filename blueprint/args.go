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
	Type        ArgType `yaml:"type"`
	Description string  `yaml:"description"`
	NoInput     bool    `yaml:"noinput"`
	Condition   string  `yaml:"condition"`
	Bool        *ArgBool
	String      *ArgString
	Array       *ArgArray
	Map         *ArgMap
}

type _Arg Arg

func (value *Arg) UnmarshalYAML(node *yaml.Node) error {
	arg := _Arg{}
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be mapping")
	}

	err := node.DecodeWith(decodeLooze, &arg)
	if err != nil {
		return err
	}

	switch arg.Type {
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
		return yamlError(node, fmt.Sprintf(`unknown argument type: %s`, arg.Type))
	}

	*value = Arg(arg)
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

type ArgString struct {
	Values  []string `yaml:"values"`
	Default *string  `yaml:"default"`
}

type ArgArray struct {
	Values  []string `yaml:"values"`
	Default []string `yaml:"default"`
}

type ArgBool struct {
	Default *bool `yaml:"default"`
}

type ArgMap struct {
	Keys Args `yaml:"keys"`
}

func NamedStringArg(name string, description string, noinput bool, condition string, values []string, defaultValue *string) NamedArg {
	return NamedArg{
		Name: name,
		Arg:  StringArg(description, noinput, condition, values, defaultValue),
	}
}

func StringArg(description string, noinput bool, condition string, values []string, defaultValue *string) Arg {
	return Arg{
		Type:        ArgTypeString,
		Description: description,
		NoInput:     noinput,
		Condition:   condition,
		String:      &ArgString{values, defaultValue},
	}
}

func NamedBooleanArg(name string, description string, noinput bool, condition string, defaultValue *bool) NamedArg {
	return NamedArg{
		Name: name,
		Arg:  BooleanArg(description, noinput, condition, defaultValue),
	}
}

func BooleanArg(description string, noinput bool, condition string, defaultValue *bool) Arg {
	return Arg{
		Type:        ArgTypeBoolean,
		Description: description,
		NoInput:     noinput,
		Condition:   condition,
		Bool:        &ArgBool{defaultValue},
	}
}

func NamedArrayArg(name string, description string, noinput bool, condition string, values []string, defaultValue []string) NamedArg {
	return NamedArg{
		Name: name,
		Arg:  ArrayArg(description, noinput, condition, values, defaultValue),
	}
}

func ArrayArg(description string, noinput bool, condition string, values []string, defaultValue []string) Arg {
	return Arg{
		Type:        ArgTypeArray,
		Description: description,
		NoInput:     noinput,
		Condition:   condition,
		Array:       &ArgArray{values, defaultValue},
	}
}

func NamedMapArg(name string, description string, noinput bool, condition string, keys Args) NamedArg {
	return NamedArg{
		Name: name,
		Arg:  MapArg(description, noinput, condition, keys),
	}
}

func MapArg(description string, noinput bool, condition string, keys Args) Arg {
	return Arg{
		Type:        ArgTypeMap,
		Description: description,
		NoInput:     noinput,
		Condition:   condition,
		Map:         &ArgMap{keys},
	}
}
