package types

import (
	"context"

	"github.com/Umaaz/redfish/pkg/format/junit"
)

type Response struct {
	Error string
	Code  int
}

type Producer interface {
	Send(ctx context.Context, results junit.TestResults) (*Response, error)
}
