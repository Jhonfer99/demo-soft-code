package infrastructure

import (
	"context"

	models "github.com/fsangopanta/demo-soft-code/modules/chat/domain/models"
)

type LocalProcessor struct{}

func (p *LocalProcessor) Process(
	ctx context.Context,
	msg models.IncomingMessage,
) (string, error) {
	return "Hola: " + msg.UserId, nil
}
