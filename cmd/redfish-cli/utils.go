package main

import (
	"context"
	"github.com/Umaaz/redfish/pkg/config"
	"os"
	"path/filepath"
)

func LoadConfig(file string, opts pklOptions) (config.App, error) {
	abs, err := filepath.Abs(opts.Pkl)
	if err != nil {
		return nil, err
	}
	_ = os.Setenv("PKL_EXEC", abs)

	ret, err := config.LoadFromPath(context.Background(), file)
	if err != nil {
		return nil, err
	}
	impl := ret.(*config.AppImpl)
	impl.File = &file
	return ret, nil
}
