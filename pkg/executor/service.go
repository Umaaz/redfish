package executor

import "github.com/Umaaz/redfish/pkg/config/pkl/gen/appconfig"

type Config struct{}

type Service struct {
	cfg Config
	app appconfig.App
}

func NewService(config Config, app appconfig.App) (*Service, error) {
	return &Service{
		cfg: config,
		app: app,
	}, nil
}

func (s *Service) Run() error {
	return nil
}
