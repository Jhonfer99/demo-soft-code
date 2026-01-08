package usecases

import (
	"context"

	domains "github.com/fsangopanta/demo-soft-code/common/domains"
)

type Processor interface {
	Process(ctx context.Context, text string, customData []domains.CustomData) (string, error)
}