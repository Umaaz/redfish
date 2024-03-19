package manager

import "github.com/Umaaz/redfish/pkg/config"

type Manager interface {
	RunJob(config *config.JobConfig) (*JobResult, error)
	RunAllJobs() ([]*JobResult, error)
}

type JobResult struct {
	Config  *config.JobConfig
	Results []*TestResult
}

type TestResult struct {
	Test       *config.Test
	Assertions []*Assertion
}

type Assertion struct {
	Pass    bool
	Name    string
	Message string
}
