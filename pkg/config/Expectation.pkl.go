// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type Expectation struct {
	Status int `pkl:"status"`

	Headers *map[string]Matcher `pkl:"headers"`

	Body *Matcher `pkl:"body"`
}
