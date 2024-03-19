// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type JobConfig struct {
	Name *string `pkl:"name" json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`

	Tests []*Test `pkl:"tests" json:"tests,omitempty" toml:"tests,omitempty" yaml:"tests,omitempty"`
}
