package junit

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/Umaaz/redfish/pkg/manager"
)

type TestResults struct {
	XMLName xml.Name `xml:"testsuites"`

	Time       int64     `xml:"time,attr"`
	Name       string    `xml:"name,omitempty,attr"`
	Tests      uint      `xml:"tests,omitempty,attr"`
	Failures   uint      `xml:"failures,omitempty,attr"`
	Errors     uint      `xml:"errors,omitempty,attr"`
	Skipped    uint      `xml:"skipped,omitempty,attr"`
	Assertions uint      `xml:"assertions,omitempty,attr"`
	Timestamp  time.Time `xml:"timestamp,omitempty,attr"`
	File       string    `xml:"file,omitempty,attr"`

	TestSuites []TestSuite
}

type TestSuite struct {
	XMLName xml.Name `xml:"testsuite"`

	Time       int64     `xml:"time,attr"`
	Name       string    `xml:"name,omitempty,attr"`
	Tests      uint      `xml:"tests,omitempty,attr"`
	Failures   uint      `xml:"failures,omitempty,attr"`
	Errors     uint      `xml:"errors,omitempty,attr"`
	Skipped    uint      `xml:"skipped,omitempty,attr"`
	Assertions uint      `xml:"assertions,omitempty,attr"`
	Timestamp  time.Time `xml:"timestamp,omitempty,attr"`
	File       string    `xml:"file,omitempty,attr"`

	TestCases []TestCase
}

type Property struct {
	XMLName xml.Name `xml:"property"`

	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr,omitempty"`
	Text  string `xml:",chardata"`
}

type TestCase struct {
	XMLName xml.Name `xml:"testcase"`

	File       string `xml:"file,attr,omitempty"`
	Name       string `xml:"name,attr,omitempty"`
	Time       int64  `xml:"time,attr,omitempty"`
	Assertions uint   `xml:"assertions,omitempty,attr"`
	Line       uint   `xml:"line,omitempty,attr"`

	FailureCount uint `xml:"-"`
	ErrorCount   uint `xml:"-"`

	Failures   []Failure
	SystemOut  SystemOut `xml:"system-out"`
	SystemErr  SystemOut `xml:"system-err"`
	Properties struct {
		Properties []Property
	} `xml:"properties"`
}

type SystemOut struct {
	Text string `xml:",cdata"`
}

type Failure struct {
	XMLName xml.Name `xml:"failure"`

	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

func Convert(results manager.TestContext) (TestResults, error) {
	totals := struct {
		tests      uint
		failures   uint
		errors     uint
		skipped    uint
		assertions uint
	}{}

	suites := make([]TestSuite, len(results.Results))

	for i, result := range results.Results {
		suite := parseResult(result)
		totals.tests += suite.Tests
		totals.failures += suite.Failures
		totals.errors += suite.Errors
		totals.skipped += suite.Skipped
		totals.assertions += suite.Assertions
		suites[i] = suite
	}

	return TestResults{
		Time:      results.Duration,
		Name:      results.Name,
		File:      results.File,
		Timestamp: results.Timestamp,

		Tests:      totals.tests,
		Failures:   totals.failures,
		Errors:     totals.errors,
		Skipped:    totals.skipped,
		Assertions: totals.assertions,

		TestSuites: suites,
	}, nil
}

func parseResult(result *manager.JobResult) TestSuite {
	totals := struct {
		tests      uint
		failures   uint
		errors     uint
		skipped    uint
		assertions uint
	}{}
	testCases := make([]TestCase, len(result.Results))
	for i, testResult := range result.Results {
		testCase, errors, skipped := parseTestResult(result, testResult)
		testCases[i] = testCase

		totals.tests += 1
		totals.failures += uint(len(testCase.Failures))
		totals.errors += errors
		totals.skipped += skipped
		totals.assertions += testCase.Assertions
	}

	return TestSuite{
		Time:       result.Duration,
		Name:       result.Name,
		Tests:      totals.tests,
		Failures:   totals.failures,
		Errors:     totals.errors,
		Skipped:    totals.skipped,
		Assertions: totals.assertions,
		Timestamp:  result.Timestamp,
		File:       result.File,
		TestCases:  testCases,
	}
}

func parseTestResult(job *manager.JobResult, result *manager.TestResult) (TestCase, uint, uint) {
	var failures []Failure
	errors := uint(0)

	for _, assertion := range result.Assertions {
		if !assertion.Pass {
			failures = append(failures, Failure{
				Message: assertion.Message,
				Type:    assertion.Type,
			})
		}
		if assertion.Error {
			errors += 1
		}
	}

	return TestCase{
		File:         job.File,
		Name:         *result.Test.Name,
		Time:         result.Duration,
		Assertions:   uint(len(result.Assertions)),
		Line:         0,
		Failures:     failures,
		ErrorCount:   errors,
		FailureCount: uint(len(failures)),
		SystemOut:    SystemOut{},
		SystemErr:    SystemOut{},
		Properties: struct{ Properties []Property }{Properties: []Property{
			{Name: "url", Value: asString(result.Test.Url)},
			{Name: "method", Value: result.Test.Method},
			{Name: "body", Text: asString(result.Test.Body)},
			{Name: "headers", Text: asString(result.Test.Headers)},
		}},
	}, errors, 0
}

func asString(value any) string {
	if value == nil {
		return ""
	}
	if value == nil {
		return ""
	}
	if v, ok := value.(string); ok {
		return v
	}
	if v, ok := value.(*string); ok {
		return *v
	}
	marshal, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%+v", value)
	}
	return string(marshal)
}
