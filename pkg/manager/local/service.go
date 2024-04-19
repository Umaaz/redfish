package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Umaaz/redfish/pkg/config/pkl/gen/appconfig"
	"github.com/Umaaz/redfish/pkg/config/pkl/gen/jobconfig"
	"github.com/Umaaz/redfish/pkg/manager"
	"github.com/Umaaz/redfish/pkg/utils/logging"
	"github.com/google/uuid"
)

var nonValueErr = errors.New("non value error")

type Config struct{}

type Service struct {
	manager.Manager

	cfg *Config
	app appconfig.App

	results []*manager.JobResult
}

func NewService(config *Config, app appconfig.App) (manager.Manager, error) {
	return &Service{
		cfg: config,
		app: app,
	}, nil
}

func (s *Service) RunJob(config *jobconfig.JobConfig, name string, i int) (*manager.JobResult, error) {
	logging.Logger.Info("Running Job", "app", s.app.GetName(), "job", config.Name)
	result := &manager.JobResult{
		Config:    config,
		Timestamp: time.Now(),
	}
	if config.Name != "" {
		result.Name = config.Name
	} else {
		config.Name = fmt.Sprintf("%s Suite: %d", name, i)
	}

	client := http.Client{
		Timeout: time.Second * 5,
	}

	for _, test := range config.Tests {
		tStart := time.Now()
		if test.Id == nil {
			u := uuid.New().String()
			test.Id = &u
		}

		reqUrl, err := s.readValue(result, test.Url)
		modifiedUrl := fmt.Sprintf("%s", reqUrl)
		if test.Name == nil {
			name := fmt.Sprintf("%s:%s", test.Method, modifiedUrl)
			test.Name = &name
		}
		logging.Logger.Info("Running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}

		body, content, err := s.processRequestBody(test, result)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}

		request, err := http.NewRequest(strings.ToUpper(test.Method), modifiedUrl, body)
		if err != nil {
			err = s.handleError(result, test, err)
			if err != nil {
				logging.Logger.Info("Error running test", "app", s.app.GetName(), "job", config.Name, "test", *test.Name, "id", *test.Id, "err", err)
				return nil, err
			}
			return result, nil
		}

		if content != "" {
			request.Header.Add("content-type", content)
		}

		if s.app.GetDefaults() != nil && s.app.GetDefaults().Headers != nil {
			for k, v := range s.app.GetDefaults().Headers {
				request.Header.Add(k, v)
			}
		}

		if test.Headers != nil {
			for key, val := range *test.Headers {
				request.Header.Add(key, val)
			}
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
		testResult.Duration = time.Now().Sub(tStart).Milliseconds()
		result.Results = append(result.Results, testResult)
	}
	result.Duration = time.Now().Sub(result.Timestamp).Milliseconds()
	logging.Logger.Info(fmt.Sprintf("test complete in (%d)ms", result.Duration))
	return result, nil
}

func (s *Service) RunAllJobs() (manager.TestContext, error) {
	start := time.Now()
	jobs := s.app.GetJobs()
	logging.Logger.Info("Running Jobs in app.", "app", s.app.GetName())
	results := make([]*manager.JobResult, len(jobs))
	for i, job := range jobs {
		res, _ := s.RunJob(job, s.app.GetName(), i)
		results[i] = res
	}

	file := s.app.GetFile()
	return manager.TestContext{
		Name:      s.app.GetName(),
		Timestamp: start,
		Duration:  time.Now().Sub(start).Milliseconds(),
		File:      *file,
		Results:   results,
	}, nil
}

func (s *Service) handleError(result *manager.JobResult, test *jobconfig.Test, err error) error {
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

func (s *Service) readValue(result *manager.JobResult, input any) (any, error) {
	if v, ok := input.(string); ok {
		return v, nil
	}
	if v, ok := input.(int); ok {
		return v, nil
	}
	if v, ok := input.(bool); ok {
		return v, nil
	}
	if v, ok := input.(float64); ok {
		return v, nil
	}

	if v, ok := input.(*jobconfig.DataSource); ok {
		id := v.SourceId
		var source *manager.TestResult
		if id == "" {
			source = result.Results[len(result.Results)-1]
		} else {
			for _, testResult := range result.Results {
				if *testResult.Test.Id == id {
					source = testResult
					break
				}
			}
		}

		if source == nil {
			return "", errors.New("cannot find datasource source: " + v.SourceId + ":" + v.Source)
		}

		switch v.Source {
		case "response":
			return s.readResponseValue(v, source)
		}
	}

	if v, ok := input.(*jobconfig.FormattingDataSource); ok {
		id := v.SourceId
		var source *manager.TestResult
		if id == "" {
			source = result.Results[len(result.Results)-1]
		} else {
			for _, testResult := range result.Results {
				if *testResult.Test.Id == id {
					source = testResult
					break
				}
			}
		}

		if source == nil {
			return "", errors.New("cannot find datasource source: " + v.SourceId + ":" + v.Source)
		}

		switch v.Source {
		case "response":
			return s.readResponseValueF(v, source)
		}
	}

	return "", nonValueErr
}

func (s *Service) checkResponse(res *http.Response, test *jobconfig.Test) (*manager.TestResult, error) {
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
	result.Assertions = append(result.Assertions, s.logAssertions(assertStatus)...)

	if test.Expected.Body != nil {
		result.Assertions = append(result.Assertions, s.logAssertions(s.checkResponseBody(res, result, test)...)...)
	}
	return result, nil
}

func (s *Service) checkResponseBody(res *http.Response, result *manager.TestResult, test *jobconfig.Test) []*manager.Assertion {
	var asserts []*manager.Assertion

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		asserts = append(asserts, &manager.Assertion{
			Pass:    false,
			Error:   false,
			Type:    "BodyReadErr",
			Message: err.Error(),
		})

		return asserts
	}
	result.Response = bodyBytes

	body := *test.Expected.Body
	if v, ok := body.(*jobconfig.JsonMatcherImpl); ok {
		var data map[string]any

		err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&data)
		if err != nil {
			asserts = append(asserts, &manager.Assertion{
				Pass:    false,
				Error:   false,
				Type:    "JSONBodyDecode",
				Message: err.Error(),
			})
		} else {
			for path, expected := range v.Expected {
				var elem any = data
				found := true
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
						if expected == nil && v == nil {
							found = true
							elem = nil
							break
						}
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

func (s *Service) processRequestBody(test *jobconfig.Test, result *manager.JobResult) (io.Reader, string, error) {
	if test.Body == nil {
		return nil, "", nil
	}

	body := *test.Body
	switch body.GetType() {
	case "json":
		return s.processJsonBody(body.(jobconfig.JsonBody), result)
	case "form":
		return s.processFormBody(body.(jobconfig.FormBody), result)
	}
	return nil, "", nil
}

func (s *Service) processJsonBody(body jobconfig.JsonBody, result *manager.JobResult) (io.Reader, string, error) {
	processBody, err := s.processObject(body.GetParams(), result)
	if err != nil {
		return nil, "", errors.Join(errors.New("failed to process json input"), err)
	}

	marshal, err := json.Marshal(processBody)
	if err != nil {
		return nil, "", errors.Join(errors.New("failed to marshal json"), err)
	}

	return bytes.NewReader(marshal), "application/json", nil
}

func (s *Service) processObject(params map[any]any, result *manager.JobResult) (map[string]any, error) {
	nMap := make(map[string]any, len(params))
	for k, v := range params {

		key, err := s.readValue(result, k)
		if err != nil {
			return nil, errors.Join(errors.New("cannot process key for object"), err)
		}

		value, err := s.readObjectValue(v, result)
		if err != nil {
			return nil, errors.Join(errors.New("cannot process value for object"), err)
		}
		nMap[fmt.Sprintf("%s", key)] = value
	}
	return nMap, nil
}

func (s *Service) readObjectValue(input any, result *manager.JobResult) (any, error) {
	value, err := s.readValue(result, input)
	if err != nil {
		if errors.Is(err, nonValueErr) {
			if val, ok := input.(map[any]any); ok {
				child, err := s.processObject(val, result)
				if err != nil {
					return nil, errors.Join(errors.New("cannot process object"), err)
				}
				return child, nil
			}
			if val, ok := input.([]any); ok {
				child, err := s.processSlice(val, result)
				if err != nil {
					return nil, errors.Join(errors.New("cannot process object"), err)
				}
				return child, nil
			}
		}
		return nil, errors.Join(errors.New("cannot process object"), err)
	}
	return value, nil
}

func (s *Service) processSlice(val []any, result *manager.JobResult) ([]any, error) {
	nSlice := make([]any, len(val))

	for i, v := range val {
		value, err := s.readObjectValue(v, result)
		if err != nil {
			return nil, errors.Join(errors.New("cannot process value for object"), err)
		}
		nSlice[i] = value
	}

	return nSlice, nil
}

func (s *Service) readResponseValue(v *jobconfig.DataSource, source *manager.TestResult) (string, error) {
	body := source.Response
	if body == nil {
		return "", errors.New("no response from source: " + *source.Test.Id)
	}

	switch v.Extractor.GetType() {
	case "json":
		var data map[string]any
		err := json.Unmarshal(body, &data)
		if err != nil {
			return "", errors.New("invalid json response from source: " + *source.Test.Id)
		}
		return s.extractFromJsonObject(v.Extractor.(jobconfig.JsonExtractor), data)
	}
	return "", nil
}

func (s *Service) extractFromJsonObject(extractor jobconfig.JsonExtractor, data map[string]any) (string, error) {
	path := extractor.GetPath()

	haystack := data
	for part := range s.walkPath(path) {
		if v, ok := haystack[part]; ok {
			if v, k := v.(string); k {
				return v, nil
			}
			if v, k := v.(map[string]any); k {
				haystack = v
				continue
			}
			if v, k := v.([]any); k {
				haystack = make(map[string]any, len(v))
				for i, item := range v {
					haystack[strconv.Itoa(i)] = item
				}
				continue
			}
		}
		return "", errors.New(fmt.Sprintf("path %s not found in object", part))
	}
	return "", errors.New(fmt.Sprintf("path %s not found in object", path))
}

func (s *Service) processFormBody(body jobconfig.FormBody, result *manager.JobResult) (io.Reader, string, error) {
	processBody, err := s.processObject(body.GetParams(), result)
	if err != nil {
		return nil, "", errors.Join(errors.New("failed to process json input"), err)
	}

	form := url.Values{}

	for k, val := range processBody {
		if v, ok := val.(string); ok {
			form.Add(k, v)
			continue
		}
		return nil, "", errors.New(fmt.Sprintf("non string value for form body %s", val))
	}

	return strings.NewReader(form.Encode()), "application/x-www-form-urlencoded", nil
}

func (s *Service) logAssertions(asserts ...*manager.Assertion) []*manager.Assertion {
	for _, assert := range asserts {
		if !assert.Pass {
			logging.Logger.Error("assertion failed", "assert", assert.Message)
		}
	}
	return asserts
}

func (s *Service) readResponseValueF(v *jobconfig.FormattingDataSource, source *manager.TestResult) (string, error) {
	body := source.Response
	if body == nil {
		return "", errors.New("no response from source: " + *source.Test.Id)
	}

	var results []any
	extractors := v.Extractors

	for _, extractor := range extractors {
		switch extractor.GetType() {
		case "json":
			var data map[string]any
			err := json.Unmarshal(body, &data)
			if err != nil {
				return "", errors.New("invalid json response from source: " + *source.Test.Id)
			}
			result, err := s.extractFromJsonObject(extractor.(jobconfig.JsonExtractor), data)
			if err != nil {
				return "", errors.Join(errors.New("cannot process json response from source"), err)
			}
			results = append(results, result)
		default:
			return "", errors.New(fmt.Sprintf("unsupported extractor type %s", extractor.GetType()))
		}
	}
	sprintf := fmt.Sprintf(v.StringFmt, results...)
	return sprintf, nil
}
