package blueprint

import (
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
)

type PathArray []string

func (arr PathArray) Matches(value string) bool {
	if arr != nil {
		for _, prefix := range arr {
			if strings.HasPrefix(value, prefix) {
				return true
			}
		}
	}
	return false
}

func (arr PathArray) Contains(value string) bool {
	if arr != nil {
		for _, prefix := range arr {
			if value == prefix {
				return true
			}
		}
	}
	return false
}

type Blueprint struct {
	Blueprint       string            `yaml:"blueprint"`
	Name            string            `yaml:"name"`
	Title           string            `yaml:"title"`
	Roots           []string          `yaml:"roots"`
	Args            Args              `yaml:"args"`
	IgnorePaths     PathArray         `yaml:"ignore"`
	ExecutablePaths PathArray         `yaml:"executables"`
	Rename          map[string]string `yaml:"rename"`
}

func Read(blueprintContent string) (*Blueprint, error) {
	blueprint := Blueprint{}
	err := yaml.Unmarshal([]byte(blueprintContent), &blueprint)
	if err != nil {
		return nil, err
	}
	return &blueprint, nil
}
