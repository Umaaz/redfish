// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type JsonExtractor interface {
	Extractor

	GetPath() string
}

var _ JsonExtractor = (*JsonExtractorImpl)(nil)

type JsonExtractorImpl struct {
	Path string `pkl:"path" json:"path,omitempty" toml:"path,omitempty" yaml:"path,omitempty"`

	Default string `pkl:"default" json:"default,omitempty" toml:"default,omitempty" yaml:"default,omitempty"`
}

func (rcv *JsonExtractorImpl) GetPath() string {
	return rcv.Path
}

func (rcv *JsonExtractorImpl) GetDefault() string {
	return rcv.Default
}
