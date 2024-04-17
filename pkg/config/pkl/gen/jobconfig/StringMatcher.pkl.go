// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package jobconfig

type StringMatcher interface {
	Matcher

	GetExpected() string
}

var _ StringMatcher = (*StringMatcherImpl)(nil)

type StringMatcherImpl struct {
	Name string `pkl:"name" json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`

	Expected string `pkl:"expected" json:"expected,omitempty" toml:"expected,omitempty" yaml:"expected,omitempty"`
}

func (rcv *StringMatcherImpl) GetName() string {
	return rcv.Name
}

func (rcv *StringMatcherImpl) GetExpected() string {
	return rcv.Expected
}
