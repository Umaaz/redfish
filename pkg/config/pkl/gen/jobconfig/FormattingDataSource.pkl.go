// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package jobconfig

type FormattingDataSource struct {
	SourceId string `pkl:"sourceId" json:"sourceId,omitempty" toml:"sourceId,omitempty" yaml:"sourceId,omitempty"`

	Source string `pkl:"source" json:"source,omitempty" toml:"source,omitempty" yaml:"source,omitempty"`

	Extractors []Extractor `pkl:"extractors" json:"extractors,omitempty" toml:"extractors,omitempty" yaml:"extractors,omitempty"`

	StringFmt string `pkl:"stringFmt" json:"stringFmt,omitempty" toml:"stringFmt,omitempty" yaml:"stringFmt,omitempty"`
}
