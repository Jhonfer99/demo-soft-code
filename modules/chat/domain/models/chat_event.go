package models

type ChatEvent struct {
	User           User            `json:"user"`
	EventTime      string          `json:"eventTime"`
	MessagePayload *MessagePayload `json:"messagePayload,omitempty"`
	Type           string          `json:"type,omitempty"`
	Message        Message         `json:"message"`
}
