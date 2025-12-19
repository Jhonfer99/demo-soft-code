package googleworkspace

import (
	"context"

	chat "cloud.google.com/go/chat/apiv1"
	chatpb "cloud.google.com/go/chat/apiv1/chatpb"
)

type chatService struct {
	client *chat.Client
}

func NewChatService(ctx context.Context) (ChatService, error) {
	c, err := chat.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &chatService{
		client: c,
	}, nil
}

func (s *chatService) SendMessage(
	ctx context.Context,
	spaceID string,
	text string,
) (*chatpb.Message, error) {

	req := &chatpb.CreateMessageRequest{
		Parent: "spaces/" + spaceID,
		Message: &chatpb.Message{
			Text: text,
		},
	}

	resp, err := s.client.CreateMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
