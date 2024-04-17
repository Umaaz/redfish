package main

import (
	"context"
	"errors"
	"github.com/Umaaz/redfish/pkg/config/pkl/gen/appconfig"
	"github.com/Umaaz/redfish/pkg/utils/logging"
	"os"
	"path/filepath"
)

func LoadConfig(file string, opts pklOptions) (appconfig.App, error) {
	if opts.Pkl != "" {
		abs, err := filepath.Abs(opts.Pkl)
		if err != nil {
			logging.Logger.Error("cannot file pkl, set --pkl or PKL_EXEC", "err", err, "pkl", opts.Pkl)
			return nil, err
		}

		lstat, err := os.Lstat(abs)
		if err != nil {
			logging.Logger.Error("cannot file pkl, set --pkl or PKL_EXEC", "err", err, "pkl", opts.Pkl)
			return nil, err
		}
		_ = lstat.Name()
		_ = os.Setenv("PKL_EXEC", abs)
	}
	if os.Getenv("PKL_EXEC") == "" {
		logging.Logger.Debug("cannot file pkl, set --pkl or PKL_EXEC")
		return nil, errors.New("pkl executable not set")
	}

	ret, err := appconfig.LoadFromPath(context.Background(), file)
	if err != nil {
		logging.Logger.Error("cannot load config", "err", err, "config_file", file)
		return nil, err
	}
	impl := ret.(*appconfig.AppImpl)
	impl.File = &file
	return ret, nil
}
