package models

type MessagePayload struct {
	Space   Space   `json:"space"`
	Message Message `json:"message"`
}