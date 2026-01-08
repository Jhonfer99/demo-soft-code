package models

type IncomingMessage struct {
	Channel string
	UserId  string
	Text    string
	ReplyTo *string
}