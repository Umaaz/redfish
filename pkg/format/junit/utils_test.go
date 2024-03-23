package junit

import (
	"encoding/xml"
	"fmt"
	"github.com/Umaaz/redfish/pkg/config"
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

func TestAsString(t *testing.T) {
	i := 1
	f := false
	s := "a value is a string"
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "basic int",
			input:    1,
			expected: "1",
		},
		{
			name:     "basic bool",
			input:    false,
			expected: "false",
		},
		{
			name:     "basic string",
			input:    "a value is a string",
			expected: "a value is a string",
		},
		{
			name:     "basic slice",
			input:    []int{1, 2, 3, 4},
			expected: "[1,2,3,4]",
		},
		{
			name:     "basic map",
			input:    map[string]int{"one": 1, "two": 2},
			expected: "{\"one\":1,\"two\":2}",
		},
		{
			name:     "pointer int",
			input:    &i,
			expected: "1",
		},
		{
			name:     "pointer bool",
			input:    &f,
			expected: "false",
		},
		{
			name:     "pointer string",
			input:    &s,
			expected: "a value is a string",
		},
		{
			name:     "pointer slice",
			input:    &[]int{1, 2, 3, 4},
			expected: "[1,2,3,4]",
		},
		{
			name:     "pointer map",
			input:    &map[string]int{"one": 1, "two": 2},
			expected: "{\"one\":1,\"two\":2}",
		},
		{
			name:     "any map",
			input:    &map[any]any{"one": 1, "two": 2},
			expected: "&map[one:1 two:2]",
		},
		{
			name:     "any map",
			input:    map[any]any{"one": 1, "two": 2},
			expected: "map[one:1 two:2]",
		},
		{
			name:     "nil",
			input:    nil,
			expected: "",
		},
		{
			name:     "nil",
			input:    config.Test{}.Body,
			expected: "null",
		},
		{
			name:     "nil",
			input:    config.Test{}.Headers,
			expected: "null",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := asString(test.input)
			if s != test.expected {
				t.Fatalf("%s != %s", s, test.expected)
			}
		})
	}
}
