package main

import (
	"encoding/json"
	"fmt"
)

type Update struct {
	pklOptions

	File string `type:"path" arg:"" help:"The path to the config file. (json, yaml, pkl)"`
}

func (u Update) Run(opts *globalOptions) error {
	loadConfig, err := LoadConfig(u.File, u.pklOptions)
	if err != nil {
		return nil
	}

	indent, err := json.MarshalIndent(loadConfig, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", string(indent))

	return nil
}
