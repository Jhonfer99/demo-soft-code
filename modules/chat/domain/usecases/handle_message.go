package usecases

import (
	"context"
	"strings"

	domains "github.com/fsangopanta/demo-soft-code/common/domains"
	models "github.com/fsangopanta/demo-soft-code/modules/chat/domain/models"
	// 	usecases "github.com/fsangopanta/demo-soft-code/modules/chat/domain/models/google/inbound"
)

type ProcessIncomingMessageUseCase struct {
	processor Processor
}

func NewProcessIncomingMessageUseCase(processor Processor) *ProcessIncomingMessageUseCase {
	return &ProcessIncomingMessageUseCase{processor: processor}
}

func (uc *ProcessIncomingMessageUseCase) Handle(
	ctx context.Context,
	msg models.IncomingMessage,
	cd []domains.CustomData,
) (models.OutgoingMessage, error) {

	text := strings.TrimSpace(msg.Text)
	if text == "" {
		return models.OutgoingMessage{Text: "Empty message"}, nil
	}
	reply, err := uc.processor.Process(ctx, msg)
	if err != nil {
		return models.OutgoingMessage{Text: "Error while processing message"}, nil
	}

	return models.OutgoingMessage{Text: reply}, nil
}
