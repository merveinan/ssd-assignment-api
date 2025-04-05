package models

import "gopkg.in/yaml.v3"

type StringSlice []string

func (s *StringSlice) UnmarshalYAML(value *yaml.Node) error {
	var multi []string
	var single string

	if err := value.Decode(&multi); err == nil {
		*s = multi
		return nil
	}

	if err := value.Decode(&single); err != nil {
		return err
	}

	*s = []string{single}
	return nil
}

type SpecificConfig struct {
	ID         string     `yaml:"id" json:"id"`
	DataSource DataSource `yaml:"datasource" json:"datasource"`
}

type DataSource struct {
	Pages map[string]StringSlice `yaml:"pages" json:"pages"`
	URLs  map[string]StringSlice `yaml:"urls" json:"urls"`
	Hosts map[string]StringSlice `yaml:"hosts" json:"hosts"`
}
