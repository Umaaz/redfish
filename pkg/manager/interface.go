package manager

import (
	"github.com/Umaaz/redfish/pkg/config"
	"time"
)

type Manager interface {
	RunJob(config *config.JobConfig) (*JobResult, error)
	RunAllJobs() (TestContext, error)
}

type TestContext struct {
	Name      string
	Timestamp time.Time
	Time      float64
	File      string

	Results []*JobResult
}

type JobResult struct {
	Name      string
	Timestamp time.Time
	Time      float64
	File      string
	Config    *config.JobConfig
	Results   []*TestResult
}

type TestResult struct {
	Test       *config.Test
	Time       float64
	Assertions []*Assertion
}

type Assertion struct {
	Pass    bool
	Error   bool
	Type    string
	Message string
}
