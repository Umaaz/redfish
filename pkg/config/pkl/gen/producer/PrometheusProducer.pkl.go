// Code generated from Pkl module `pickle.config.producer`. DO NOT EDIT.
package producer

type PrometheusProducer interface {
	HttpEndpoint

	GetType() any

	GetLabels() []*Label

	GetCustomLables() *[]*Label
}

var _ PrometheusProducer = (*PrometheusProducerImpl)(nil)

type PrometheusProducerImpl struct {
	Type any `pkl:"type" json:"type,omitempty" toml:"type,omitempty" yaml:"type,omitempty"`

	Path string `pkl:"path" json:"path,omitempty" toml:"path,omitempty" yaml:"path,omitempty"`

	Labels []*Label `pkl:"labels" json:"labels,omitempty" toml:"labels,omitempty" yaml:"labels,omitempty"`

	CustomLables *[]*Label `pkl:"customLables" json:"customLables,omitempty" toml:"customLables,omitempty" yaml:"customLables,omitempty"`

	Url string `pkl:"url" json:"url,omitempty" toml:"url,omitempty" yaml:"url,omitempty"`

	Timeout int `pkl:"timeout" json:"timeout,omitempty" toml:"timeout,omitempty" yaml:"timeout,omitempty"`

	Headers *map[string]string `pkl:"headers" json:"headers,omitempty" toml:"headers,omitempty" yaml:"headers,omitempty"`
}

func (rcv *PrometheusProducerImpl) GetType() any {
	return rcv.Type
}

func (rcv *PrometheusProducerImpl) GetPath() string {
	return rcv.Path
}

func (rcv *PrometheusProducerImpl) GetLabels() []*Label {
	return rcv.Labels
}

func (rcv *PrometheusProducerImpl) GetCustomLables() *[]*Label {
	return rcv.CustomLables
}

func (rcv *PrometheusProducerImpl) GetUrl() string {
	return rcv.Url
}

func (rcv *PrometheusProducerImpl) GetTimeout() int {
	return rcv.Timeout
}

func (rcv *PrometheusProducerImpl) GetHeaders() *map[string]string {
	return rcv.Headers
}
