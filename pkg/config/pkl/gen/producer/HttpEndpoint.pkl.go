// Code generated from Pkl module `pickle.config.producer`. DO NOT EDIT.
package producer

type HttpEndpoint interface {
	GetUrl() string

	GetPath() string

	GetTimeout() int

	GetHeaders() *map[string]string
}
