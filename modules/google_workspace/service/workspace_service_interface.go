package googleworkspace

import (
	"context"

	chatpb "cloud.google.com/go/chat/apiv1/chatpb"
)

type ChatService interface {
	SendMessage(ctx context.Context, spaceID string, text string) (*chatpb.Message, error)
}
