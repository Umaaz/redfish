// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package jobconfig

type FormBody interface {
	RequestBody

	GetParams() map[any]any
}

var _ FormBody = (*FormBodyImpl)(nil)

type FormBodyImpl struct {
	Type string `pkl:"type" json:"type,omitempty" toml:"type,omitempty" yaml:"type,omitempty"`

	Params map[any]any `pkl:"params" json:"params,omitempty" toml:"params,omitempty" yaml:"params,omitempty"`
}

func (rcv *FormBodyImpl) GetType() string {
	return rcv.Type
}

func (rcv *FormBodyImpl) GetParams() map[any]any {
	return rcv.Params
}
