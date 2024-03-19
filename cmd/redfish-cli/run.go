package main

import (
	"encoding/json"
	"fmt"
	"github.com/Umaaz/redfish/pkg/manager/local"
)

type Run struct {
	pklOptions

	File string `type:"path" arg:"" help:"The path to the config file. (json, yaml, pkl)"`
}

func (r Run) Run(opts *globalOptions) error {

	loadConfig, err := LoadConfig(r.File, r.pklOptions)
	if err != nil {
		return nil
	}

	service, err := local.NewService(nil, loadConfig)
	if err != nil {
		return err
	}

	results, err := service.RunAllJobs()
	if err != nil {
		return err
	}

	indent, err := json.MarshalIndent(results, "", "  ")
	fmt.Printf("%+v\n", string(indent))

	return nil
}
