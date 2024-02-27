// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type StringMatcher interface {
	Matcher

	GetExpected() string
}

var _ StringMatcher = (*StringMatcherImpl)(nil)

type StringMatcherImpl struct {
	Name string `pkl:"name"`

	Expected string `pkl:"expected"`
}

func (rcv *StringMatcherImpl) GetName() string {
	return rcv.Name
}

func (rcv *StringMatcherImpl) GetExpected() string {
	return rcv.Expected
}
