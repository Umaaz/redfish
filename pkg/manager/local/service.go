package local

import (
	"fmt"
	"github.com/Umaaz/redfish/pkg/config"
	"github.com/Umaaz/redfish/pkg/manager"
	"github.com/Umaaz/redfish/pkg/utils/logging"
	"github.com/google/uuid"
	"net/http"
	"time"
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
	logging.Logger.Info("Running Job", "app", s.app.GetName(), "job", config.Name)
	result := &manager.JobResult{
		Config: config,
	}

	for _, test := range config.Tests {
		if test.Id == nil {
			u := uuid.New().String()
			test.Id = &u
		}
		if test.Name == nil {
			name := fmt.Sprintf("%s:%s", test.Method, test.Url)
			test.Name = &name
		}
		logging.Logger.Info("Running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id)

		client := http.Client{
			Timeout: time.Second * 5,
		}
		request, err := http.NewRequest(test.Method, s.readValue(result, test.Url), nil)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}

		res, err := client.Do(request)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}

		testResult, err := s.checkResponse(res, test)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}
		result.Results = append(result.Results, testResult)
	}

	return result, nil
}

func (s *Service) RunAllJobs() (manager.TestContext, error) {
	start := time.Now()
	jobs := s.app.GetJobs()
	logging.Logger.Info("Running Jobs in app.", "app", s.app.GetName())
	results := make([]*manager.JobResult, len(jobs))
	for i, job := range jobs {
		res, _ := s.RunJob(job)
		results[i] = res
	}

	file := s.app.GetFile()
	return manager.TestContext{
		Name:      s.app.GetName(),
		Timestamp: start,
		Time:      time.Now().Sub(start).Seconds(),
		File:      *file,
		Results:   results,
	}, nil
}

func (s *Service) handleError(result *manager.JobResult, test *config.Test, err error) error {
	logging.Logger.Info("handling error from test", "app", s.app.GetName(), "job", result.Config.Name, "test", *test.Name, "id", *test.Id, "err", err)
	tResult := &manager.TestResult{
		Test: test,
		Assertions: []*manager.Assertion{{
			Pass:    false,
			Error:   true,
			Type:    "ExecutingError",
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
		Type:    "StatusCode",
		Message: fmt.Sprintf("%d == %d", test.Expected.Status, res.StatusCode),
	}
	if res.StatusCode != test.Expected.Status {
		assertStatus.Message = fmt.Sprintf("%d != %d", test.Expected.Status, res.StatusCode)
	}
	result.Assertions = append(result.Assertions, assertStatus)

	return result, nil
}
