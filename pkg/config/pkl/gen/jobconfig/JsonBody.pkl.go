// Code generated from Pkl module `pickle.config.job`. DO NOT EDIT.
package jobconfig

type JsonBody interface {
	RequestBody

	GetParams() map[any]any
}

var _ JsonBody = (*JsonBodyImpl)(nil)

type JsonBodyImpl struct {
	Type string `pkl:"type" json:"type,omitempty" toml:"type,omitempty" yaml:"type,omitempty"`

	Params map[any]any `pkl:"params" json:"params,omitempty" toml:"params,omitempty" yaml:"params,omitempty"`
}

func (rcv *JsonBodyImpl) GetType() string {
	return rcv.Type
}

func (rcv *JsonBodyImpl) GetParams() map[any]any {
	return rcv.Params
}
