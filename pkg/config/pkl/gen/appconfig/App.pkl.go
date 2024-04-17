// Code generated from Pkl module `redfish.config.app`. DO NOT EDIT.
package appconfig

import (
	"context"

	"github.com/Umaaz/redfish/pkg/config/pkl/gen/jobconfig"
	"github.com/Umaaz/redfish/pkg/config/pkl/gen/producer"
	"github.com/apple/pkl-go/pkl"
)

type App interface {
	jobconfig.Job

	GetFile() *string

	GetName() string

	GetJobs() []*jobconfig.JobConfig

	GetRules() []*jobconfig.RuleConfig

	GetProducer() *producer.Producer
}

var _ App = (*AppImpl)(nil)

type AppImpl struct {
	File *string `pkl:"file" json:"file,omitempty" toml:"file,omitempty" yaml:"file,omitempty"`

	Name string `pkl:"name" json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`

	Jobs []*jobconfig.JobConfig `pkl:"jobs" json:"jobs,omitempty" toml:"jobs,omitempty" yaml:"jobs,omitempty"`

	Rules []*jobconfig.RuleConfig `pkl:"rules" json:"rules,omitempty" toml:"rules,omitempty" yaml:"rules,omitempty"`

	Producer *producer.Producer `pkl:"producer" json:"producer,omitempty" toml:"producer,omitempty" yaml:"producer,omitempty"`
}

func (rcv *AppImpl) GetFile() *string {
	return rcv.File
}

func (rcv *AppImpl) GetName() string {
	return rcv.Name
}

func (rcv *AppImpl) GetJobs() []*jobconfig.JobConfig {
	return rcv.Jobs
}

func (rcv *AppImpl) GetRules() []*jobconfig.RuleConfig {
	return rcv.Rules
}

func (rcv *AppImpl) GetProducer() *producer.Producer {
	return rcv.Producer
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
