// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package config

type JsonMatcher interface {
	Matcher

	GetExpected() map[any]any
}

var _ JsonMatcher = (*JsonMatcherImpl)(nil)

type JsonMatcherImpl struct {
	Name string `pkl:"name"`

	Expected map[any]any `pkl:"expected"`
}

func (rcv *JsonMatcherImpl) GetName() string {
	return rcv.Name
}

func (rcv *JsonMatcherImpl) GetExpected() map[any]any {
	return rcv.Expected
}
