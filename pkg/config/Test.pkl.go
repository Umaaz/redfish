// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type Test struct {
	Id *string `pkl:"id"`

	Name *string `pkl:"name"`

	Url any `pkl:"url"`

	Method string `pkl:"method"`

	Headers *map[any]any `pkl:"headers"`

	Expected *Expectation `pkl:"expected"`
}
