package main

import (
	"encoding/xml"
	"fmt"
	"github.com/Umaaz/redfish/pkg/format/junit"
	"github.com/Umaaz/redfish/pkg/manager/local"
)

type Run struct {
	pklOptions

	File string `type:"path" arg:"" help:"The path to the config file. (json, yaml, pkl)"`
}

func (r Run) Run(opts *globalOptions) error {

	loadConfig, err := LoadConfig(r.File, r.pklOptions)
	if err != nil {
		return err
	}

	service, err := local.NewService(nil, loadConfig)
	if err != nil {
		return err
	}

	results, err := service.RunAllJobs()
	if err != nil {
		return err
	}

	convert, err := junit.Convert(results)
	if err != nil {
		return err
	}

	out, _ := xml.MarshalIndent(convert, " ", "  ")

	fmt.Println(xml.Header + string(out))

	return nil
}
