package main

import (
	"context"
	"encoding/xml"
	"os"

	"github.com/Umaaz/redfish/pkg/format/junit"
	"github.com/Umaaz/redfish/pkg/manager/local"
	"github.com/Umaaz/redfish/pkg/producer"
	"github.com/Umaaz/redfish/pkg/utils/logging"
)

type Run struct {
	pklOptions

	File   string `type:"path" arg:"" help:"The path to the config file. (json, yaml, pkl)"`
	Output string `type:"path" help:"The path to the output file" default:"-"`
}

func (r Run) Run(opts *globalOptions) error {
	loadConfig, err := LoadConfig(r.File, r.pklOptions)
	if err != nil {
		logging.Logger.Error("cannot load config", "err", err)
		return err
	}

	service, err := local.NewService(nil, loadConfig)
	if err != nil {
		logging.Logger.Error("cannot load service", "err", err)
		return err
	}

	results, err := service.RunAllJobs()
	if err != nil {
		logging.Logger.Error("cannot run jobs", "err", err)
		return err
	}

	convert, err := junit.Convert(results)
	if err != nil {
		logging.Logger.Error("cannot convert to junit", "err", err)
		return err
	}

	if r.Output != "-" {

		out, _ := xml.MarshalIndent(convert, " ", "  ")

		err = os.WriteFile(r.Output, []byte(xml.Header+string(out)), 0o644)
		if err != nil {
			logging.Logger.Error("cannot write file", "err", err)
			return err
		}
	}

	if loadConfig.GetProducer() != nil {
		err := producer.Run(context.Background(), loadConfig.GetProducer(), convert)
		if err != nil {
			logging.Logger.Error("failed to run producer", "err", err)
			return err
		}
	}

	return nil
}
