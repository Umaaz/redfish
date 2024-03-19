// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type Expectation struct {
	Status int `pkl:"status" json:"status,omitempty" toml:"status,omitempty" yaml:"status,omitempty"`

	Headers *map[string]Matcher `pkl:"headers" json:"headers,omitempty" toml:"headers,omitempty" yaml:"headers,omitempty"`

	Body *Matcher `pkl:"body" json:"body,omitempty" toml:"body,omitempty" yaml:"body,omitempty"`
}
