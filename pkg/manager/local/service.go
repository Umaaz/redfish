package local

import (
	"fmt"
	"github.com/Umaaz/redfish/pkg/config"
	"github.com/Umaaz/redfish/pkg/manager"
	"github.com/google/uuid"
	"net/http"
)

type Config struct {
}

type Service struct {
	manager.Manager

	cfg *Config
	app config.App

	results []*manager.JobResult
}

func NewService(config *Config, app config.App) (manager.Manager, error) {
	return &Service{
		cfg: config,
		app: app,
	}, nil
}

func (s *Service) RunJob(config *config.JobConfig) (*manager.JobResult, error) {
	result := &manager.JobResult{
		Config: config,
	}

	for _, test := range config.Tests {
		if test.Id == nil {
			u := uuid.New().String()
			test.Id = &u
		}

		client := http.Client{}
		request, err := http.NewRequest(test.Method, s.readValue(result, test.Url), nil)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		res, err := client.Do(request)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		testResult, err := s.checkResponse(res, test)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		result.Results = append(result.Results, testResult)
	}

	return result, nil
}

func (s *Service) RunAllJobs() ([]*manager.JobResult, error) {
	jobs := s.app.GetJobs()
	results := make([]*manager.JobResult, len(jobs))
	for i, job := range jobs {
		res, _ := s.RunJob(job)
		results[i] = res
	}
	return results, nil
}

func (s *Service) handleError(result *manager.JobResult, test *config.Test, err error) error {
	tResult := &manager.TestResult{
		Test: test,
		Assertions: []*manager.Assertion{{
			Pass:    false,
			Name:    "Error Executing test",
			Message: err.Error(),
		}},
	}

	result.Results = append(result.Results, tResult)

	return nil
}

func (s *Service) readValue(result *manager.JobResult, url any) string {
	if v, ok := url.(string); ok {
		return v
	}

	//todo read as datasource and load from results
	return ""
}

func (s *Service) checkResponse(res *http.Response, test *config.Test) (*manager.TestResult, error) {
	result := &manager.TestResult{
		Test: test,
	}

	assertStatus := &manager.Assertion{
		Pass:    true,
		Name:    "StatusCode",
		Message: fmt.Sprintf("%d == %d", test.Expected.Status, res.StatusCode),
	}
	if res.StatusCode != test.Expected.Status {
		assertStatus.Message = fmt.Sprintf("%d != %d", test.Expected.Status, res.StatusCode)
	}
	result.Assertions = append(result.Assertions, assertStatus)

	return result, nil
}
