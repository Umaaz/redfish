// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type JsonExtractor interface {
	Extractor

	GetPath() string
}

var _ JsonExtractor = (*JsonExtractorImpl)(nil)

type JsonExtractorImpl struct {
	Type string `pkl:"type" json:"type,omitempty" toml:"type,omitempty" yaml:"type,omitempty"`

	Path string `pkl:"path" json:"path,omitempty" toml:"path,omitempty" yaml:"path,omitempty"`

	Default string `pkl:"default" json:"default,omitempty" toml:"default,omitempty" yaml:"default,omitempty"`
}

func (rcv *JsonExtractorImpl) GetType() string {
	return rcv.Type
}

func (rcv *JsonExtractorImpl) GetPath() string {
	return rcv.Path
}

func (rcv *JsonExtractorImpl) GetDefault() string {
	return rcv.Default
}
