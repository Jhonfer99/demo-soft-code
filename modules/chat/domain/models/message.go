package models

type Message struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	Thread Thread `json:"thread"`
}
