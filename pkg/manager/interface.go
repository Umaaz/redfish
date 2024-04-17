package manager

import (
	"github.com/Umaaz/redfish/pkg/config/pkl/gen/jobconfig"
	"time"
)

type Manager interface {
	RunJob(config *jobconfig.JobConfig, name string, i int) (*JobResult, error)
	RunAllJobs() (TestContext, error)
}

type TestContext struct {
	Name      string
	Timestamp time.Time
	Duration  int64
	File      string

	Results []*JobResult
}

type JobResult struct {
	Name      string
	Timestamp time.Time
	Duration  int64
	File      string
	Config    *jobconfig.JobConfig
	Results   []*TestResult
}

type TestResult struct {
	Test       *jobconfig.Test
	Duration   int64
	Assertions []*Assertion

	Response []byte
}

type Assertion struct {
	Pass    bool
	Error   bool
	Type    string
	Message string
}
