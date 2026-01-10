package usecases

import (
	"context"

	models "github.com/fsangopanta/demo-soft-code/modules/chat/domain/models"
)

type Processor interface {
	Process(ctx context.Context, req models.IncomingMessage) (string, error)
}
