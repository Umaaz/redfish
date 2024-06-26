// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package jobconfig

type JsonMatcher interface {
	Matcher

	GetExpected() map[string]any
}

var _ JsonMatcher = (*JsonMatcherImpl)(nil)

type JsonMatcherImpl struct {
	Name string `pkl:"name" json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`

	Expected map[string]any `pkl:"expected" json:"expected,omitempty" toml:"expected,omitempty" yaml:"expected,omitempty"`
}

func (rcv *JsonMatcherImpl) GetName() string {
	return rcv.Name
}

func (rcv *JsonMatcherImpl) GetExpected() map[string]any {
	return rcv.Expected
}
