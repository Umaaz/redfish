package prometheus

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Umaaz/redfish/pkg/config/pkl/gen/producer"
	"github.com/Umaaz/redfish/pkg/format/junit"
	"github.com/Umaaz/redfish/pkg/producer/types"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	"io"
	"net/http"
	"slices"
	"time"
)

type prometheusService struct {
	producer producer.PrometheusProducer
	client   http.Client
}

var _ types.Producer = (*prometheusService)(nil)

func NewPrometheus(producer producer.PrometheusProducer) (types.Producer, error) {
	return &prometheusService{
		producer: producer,
		client:   http.Client{Timeout: time.Duration(producer.GetTimeout()) * time.Second},
	}, nil
}

func (p *prometheusService) Send(ctx context.Context, results junit.TestResults) (*types.Response, error) {
	req := &prompb.WriteRequest{
		Timeseries: make([]prompb.TimeSeries, 0),
		Metadata:   make([]prompb.MetricMetadata, 0),
	}
	for _, suite := range results.TestSuites {
		for _, testResult := range suite.TestCases {
			labels := processLabels(p.producer.GetLabels(), p.producer.GetCustomLables(), results, suite, testResult)
			req.Timeseries = append(req.Timeseries, prompb.TimeSeries{
				Labels: append(slices.Clone(labels), prompb.Label{
					Name:  "__name__",
					Value: "redfish_duration",
				}),
				Samples: []prompb.Sample{
					{
						Value:     float64(testResult.Time),
						Timestamp: results.Timestamp.UnixMilli(),
					},
				},
			})
			req.Metadata = append(req.Metadata, prompb.MetricMetadata{
				MetricFamilyName: "redfish_duration",
				Type:             prompb.MetricMetadata_GAUGE,
				Help:             "Duration of test.",
				Unit:             "ms",
			})
			req.Timeseries = append(req.Timeseries, prompb.TimeSeries{
				Labels: append(slices.Clone(labels), prompb.Label{
					Name:  "__name__",
					Value: "redfish_assertions_count",
				}),
				Samples: []prompb.Sample{
					{
						Value:     float64(testResult.Assertions),
						Timestamp: results.Timestamp.UnixMilli(),
					},
				},
			})
			req.Metadata = append(req.Metadata, prompb.MetricMetadata{
				MetricFamilyName: "redfish_assertions_count",
				Type:             prompb.MetricMetadata_GAUGE,
				Help:             "Number of assertions in test.",
				Unit:             "",
			})
			req.Timeseries = append(req.Timeseries, prompb.TimeSeries{
				Labels: append(slices.Clone(labels), prompb.Label{
					Name:  "__name__",
					Value: "redfish_assertions_failures",
				}),
				Samples: []prompb.Sample{
					{
						Value:     float64(len(testResult.Failures)),
						Timestamp: results.Timestamp.UnixMilli(),
					},
				},
			})
			req.Metadata = append(req.Metadata, prompb.MetricMetadata{
				MetricFamilyName: "redfish_assertions_failures",
				Type:             prompb.MetricMetadata_GAUGE,
				Help:             "Number of assertion failures in test.",
				Unit:             "",
			})
			req.Timeseries = append(req.Timeseries, prompb.TimeSeries{
				Labels: append(slices.Clone(labels), prompb.Label{
					Name:  "__name__",
					Value: "redfish_assertions_errors",
				}),
				Samples: []prompb.Sample{
					{
						Value:     float64(testResult.ErrorCount),
						Timestamp: results.Timestamp.UnixMilli(),
					},
				},
			})
			req.Metadata = append(req.Metadata, prompb.MetricMetadata{
				MetricFamilyName: "redfish_assertions_errors",
				Type:             prompb.MetricMetadata_GAUGE,
				Help:             "Number of assertion errors in test.",
				Unit:             "",
			})
		}
	}

	return p.write(ctx, req)
}

func (p *prometheusService) write(ctx context.Context, req *prompb.WriteRequest) (*types.Response, error) {
	// Marshal proto and compress.
	pbBytes, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("promwrite: marshaling remote write request proto: %w", err)
	}

	compressedBytes := snappy.Encode(nil, pbBytes)
	// Prepare http request.
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", p.producer.GetUrl(), p.producer.GetPath()), bytes.NewBuffer(compressedBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("User-Agent", "RedFish 0.0.0")
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	if p.producer.GetHeaders() != nil {
		for k, v := range *p.producer.GetHeaders() {
			httpReq.Header.Add(k, v)
		}
	}

	httpResp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("promwrite: sending remote write request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResp.Body)

	if st := httpResp.StatusCode; st/100 != 2 {
		msg, _ := io.ReadAll(httpResp.Body)
		return &types.Response{
			Error: fmt.Sprintf("promwrite: expected status %d, got %d: %s", http.StatusOK, st, string(msg)),
			Code:  st,
		}, nil
	}
	return &types.Response{
		Code: httpResp.StatusCode,
	}, nil
}

func processLabels(labels []*producer.Label, userLabels *[]*producer.Label, testContext junit.TestResults, suite junit.TestSuite, testResult junit.TestCase) []prompb.Label {
	var ret []prompb.Label

	for _, label := range labels {
		switch label.Name {
		case "name":
			ret = append(ret, prompb.Label{
				Name:  "name",
				Value: testContext.Name,
			})
		case "suite_name":
			ret = append(ret, prompb.Label{
				Name:  "suite_name",
				Value: suite.Name,
			})
		case "test_name":
			ret = append(ret, prompb.Label{
				Name:  "test_name",
				Value: testResult.Name,
			})
		}
	}

	if userLabels != nil {
		for _, label := range *userLabels {
			ret = append(ret, prompb.Label{
				Name:  label.Name,
				Value: label.Value,
			})
		}
	}

	return ret
}
