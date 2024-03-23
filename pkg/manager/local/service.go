package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Umaaz/redfish/pkg/config"
	"github.com/Umaaz/redfish/pkg/manager"
	"github.com/Umaaz/redfish/pkg/utils/logging"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		Config:    config,
		Timestamp: time.Now(),
	}

	client := http.Client{
		Timeout: time.Second * 5,
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

		reqUrl, err := s.readValue(result, test.Url)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}
		request, err := http.NewRequest(strings.ToUpper(test.Method), reqUrl, nil)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}

		if test.Headers != nil {
			for key, val := range *test.Headers {
				request.Header.Add(key, val)
			}
		}
		urlP, err := url.Parse(reqUrl)
		request.Host = urlP.Host

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
	result.Time = time.Now().Sub(result.Timestamp).Seconds()
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

func (s *Service) readValue(result *manager.JobResult, url any) (string, error) {
	if v, ok := url.(string); ok {
		return v, nil
	}

	if v, ok := url.(config.DataSource); ok {
		id := v.SourceId
		var source *manager.TestResult
		if id == "" {
			source = result.Results[len(result.Results)-1]
		} else {
			for _, testResult := range result.Results {
				if testResult.Test.Id == &id {
					source = testResult
					break
				}
			}
		}
		if source == nil {
			return "", errors.New("cannot find datasource source: " + v.SourceId + ":" + v.Source)
		}
	}
	return "", nil
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
		assertStatus.Pass = false
	}
	result.Assertions = append(result.Assertions, assertStatus)

	if test.Expected.Body != nil {
		result.Assertions = append(result.Assertions, s.checkResponseBody(res, test)...)
	}
	return result, nil
}

func (s *Service) checkResponseBody(res *http.Response, test *config.Test) []*manager.Assertion {
	var asserts []*manager.Assertion
	body := *test.Expected.Body
	if v, ok := body.(*config.JsonMatcherImpl); ok {
		var data map[string]any
		bodyBytes, _ := io.ReadAll(res.Body)

		err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&data)
		if err != nil {
			asserts = append(asserts, &manager.Assertion{
				Pass:    false,
				Error:   false,
				Type:    "JSONBodyDecode",
				Message: err.Error(),
			})

		} else {

			var elem any = data
			found := true
			for path, expected := range v.Expected {
				for part := range s.walkPath(path) {
					if elem == nil {
						asserts = append(asserts, &manager.Assertion{
							Pass:    false,
							Error:   false,
							Type:    "JSONBodyDecode",
							Message: fmt.Sprintf("cannot read next %s in path %s", part, path),
						})
						found = false
						break
					}
					if v, ok := elem.(map[string]any)[part]; ok {
						if v == nil {
							asserts = append(asserts, &manager.Assertion{
								Pass:    false,
								Error:   false,
								Type:    "JSONBodyDecode",
								Message: fmt.Sprintf("cannot read next %s in path %s", part, path),
							})
						}
						elem = v
					} else {
						asserts = append(asserts, &manager.Assertion{
							Pass:    false,
							Error:   false,
							Type:    "JSONBodyDecode",
							Message: fmt.Sprintf("cannot read next %s in path %s", part, path),
						})
						found = false
						break
					}
				}
				if found {
					if elem != expected {
						asserts = append(asserts, &manager.Assertion{
							Pass:    false,
							Error:   false,
							Type:    "JSONBodyDecode",
							Message: fmt.Sprintf("%s != %s at path %s", expected, elem, path),
						})
					} else {
						asserts = append(asserts, &manager.Assertion{
							Pass:    true,
							Error:   false,
							Type:    "JSONBodyDecode",
							Message: fmt.Sprintf("%s == %s at path %s", expected, elem, path),
						})
					}
				}
			}
		}
	}

	return asserts
}

func (s *Service) walkPath(path string) chan string {
	channel := make(chan string)

	// if we have no cases of \. in the path then use basic split
	if !strings.Contains(path, "\\.") {
		go func() {
			defer close(channel)
			split := strings.Split(path, ".")
			for _, s := range split {
				channel <- s
			}
		}()
		return channel
	}

	go func(str string) {
		defer close(channel)

		split := strings.Split(str, ".")
		part := ""
		for i := 0; i < len(split); i++ {
			part += split[i]
			if !strings.HasSuffix(part, "\\") {
				channel <- part
				part = ""
				continue
			}
			part += "."
		}
	}(path)

	return channel
}
