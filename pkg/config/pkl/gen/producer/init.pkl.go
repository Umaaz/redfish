// Code generated from Pkl module `pickle.config.producer`. DO NOT EDIT.
package producer

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("pickle.config.producer", Producer{})
	pkl.RegisterMapping("pickle.config.producer#PrometheusProducer", PrometheusProducerImpl{})
	pkl.RegisterMapping("pickle.config.producer#Label", Label{})
}
