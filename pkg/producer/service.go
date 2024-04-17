package producer

import (
	"context"
	"errors"
	producercfg "github.com/Umaaz/redfish/pkg/config/pkl/gen/producer"
	"github.com/Umaaz/redfish/pkg/format/junit"
	"github.com/Umaaz/redfish/pkg/producer/prometheus"
	"github.com/Umaaz/redfish/pkg/utils/logging"
)

func Run(ctx context.Context, producer *producercfg.Producer, results junit.TestResults) error {
	if producer.Prometheus != nil {
		prometheusProducer := *producer.Prometheus
		newPrometheus, err := prometheus.NewPrometheus(prometheusProducer.(producercfg.PrometheusProducer))
		if err != nil {
			return err
		}
		send, err := newPrometheus.Send(ctx, results)
		if err != nil {
			return err
		}
		if send.Error != "" {
			logging.Logger.Error(send.Error)
			return errors.New(send.Error)
		}
		return nil
	}
	return errors.New("no producer provided")
}
