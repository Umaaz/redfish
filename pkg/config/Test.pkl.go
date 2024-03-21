// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type Test struct {
	Id *string `pkl:"id" json:"id,omitempty" toml:"id,omitempty" yaml:"id,omitempty"`

	Name *string `pkl:"name" json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`

	Url any `pkl:"url" json:"url,omitempty" toml:"url,omitempty" yaml:"url,omitempty"`

	Method string `pkl:"method" json:"method,omitempty" toml:"method,omitempty" yaml:"method,omitempty"`

	Body *RequestBody `pkl:"body" json:"body,omitempty" toml:"body,omitempty" yaml:"body,omitempty"`

	Headers *map[any]any `pkl:"headers" json:"headers,omitempty" toml:"headers,omitempty" yaml:"headers,omitempty"`

	Expected *Expectation `pkl:"expected" json:"expected,omitempty" toml:"expected,omitempty" yaml:"expected,omitempty"`
}
