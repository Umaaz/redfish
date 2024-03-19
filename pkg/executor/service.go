package executor

import "github.com/Umaaz/redfish/pkg/config"

type Config struct {
}

type Service struct {
	cfg Config
	app config.App
}

func NewService(config Config, app config.App) (*Service, error) {
	return &Service{
		cfg: config,
		app: app,
	}, nil
}

func (s *Service) Run() error {
	return nil
}
