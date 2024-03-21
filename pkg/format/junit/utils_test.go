package junit

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestXmlFormat(t *testing.T) {
	results := TestResults{
		Time: 01,
		TestSuites: []TestSuite{
			{
				Name: "test1", Time: 12,
				TestCases: []TestCase{{
					Name: "testcase1",
					Time: 10,
					Failures: []Failure{{
						Message: "200 != 201",
						Type:    "StatusCodeValidation",
					}},
					SystemOut: SystemOut{
						Text: `
"this is a log message that " +
						"would be a system out or log"`,
					},
					Properties: struct{ Properties []Property }{Properties: []Property{
						{
							Name:  "request_url",
							Value: "https://example.com",
							Text:  "asldasljd",
						},
					}},
				}},
			},
		},
	}

	out, _ := xml.MarshalIndent(results, " ", "  ")

	fmt.Println(xml.Header + string(out))
}
