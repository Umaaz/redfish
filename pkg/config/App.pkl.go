// Code generated from Pkl module `redfish.config.app`. DO NOT EDIT.
package config

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type App interface {
	Job

	GetName() string

	GetJobs() []*JobConfig

	GetRules() []*RuleConfig
}

var _ App = (*AppImpl)(nil)

type AppImpl struct {
	Name string `pkl:"name"`

	Jobs []*JobConfig `pkl:"jobs"`

	Rules []*RuleConfig `pkl:"rules"`
}

func (rcv *AppImpl) GetName() string {
	return rcv.Name
}

func (rcv *AppImpl) GetJobs() []*JobConfig {
	return rcv.Jobs
}

func (rcv *AppImpl) GetRules() []*RuleConfig {
	return rcv.Rules
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a App
func LoadFromPath(ctx context.Context, path string) (ret App, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a App
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (App, error) {
	var ret AppImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
