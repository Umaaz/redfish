// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type DataSource struct {
	SourceId string `pkl:"sourceId" json:"sourceId,omitempty" toml:"sourceId,omitempty" yaml:"sourceId,omitempty"`

	Source string `pkl:"source" json:"source,omitempty" toml:"source,omitempty" yaml:"source,omitempty"`

	Extractor Extractor `pkl:"extractor" json:"extractor,omitempty" toml:"extractor,omitempty" yaml:"extractor,omitempty"`
}
