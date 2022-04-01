package blueprint

import (
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
)

type PathPrefixArray []string

func (arr PathPrefixArray) Matches(value string) bool {
	if arr != nil {
		for _, prefix := range arr {
			if strings.HasPrefix(value, prefix) {
				return true
			}
		}
	}
	return false
}

type Blueprint struct {
	Blueprint   string          `yaml:"blueprint"`
	Name        string          `yaml:"name"`
	Title       string          `yaml:"title"`
	Roots       []string        `yaml:"roots"`
	Args        Args            `yaml:"args"`
	IgnorePaths PathPrefixArray `yaml:"ignore"`
}

func Read(blueprintContent string) (*Blueprint, error) {
	blueprint := Blueprint{}
	err := yaml.Unmarshal([]byte(blueprintContent), &blueprint)
	if err != nil {
		return nil, err
	}
	return &blueprint, nil
}
