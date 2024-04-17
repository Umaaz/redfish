// Code generated from Pkl module `pickle.config.producer`. DO NOT EDIT.
package producer

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Producer struct {
	Prometheus *PrometheusProducer `pkl:"prometheus" json:"prometheus,omitempty" toml:"prometheus,omitempty" yaml:"prometheus,omitempty"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Producer
func LoadFromPath(ctx context.Context, path string) (ret *Producer, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Producer
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Producer, error) {
	var ret Producer
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
