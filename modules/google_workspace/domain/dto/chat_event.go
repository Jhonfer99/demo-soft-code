package dto

type WorkspaceEvent struct {
	Common CommonEventObject `json:"commonEventObject"`
	Chat   *ChatEvent        `json:"chat"`
}

type CommonEventObject struct {
	UserLocale string `json:"userLocale"`
	HostApp    string `json:"hostApp"`
	Platform   string `json:"platform"`
}

type ChatEvent struct {
	User           User            `json:"user"`
	EventTime      string          `json:"eventTime"`
	MessagePayload *MessagePayload `json:"messagePayload,omitempty"`
}

type MessagePayload struct {
	Space   Space   `json:"space"`
	Message Message `json:"message"`
}

type Message struct {
	Name   string `json:"name"`
	Text   string `json:"text"`
	Thread Thread `json:"thread"`
}

type Thread struct {
	Name string `json:"name"`
}

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Type        string `json:"type"`
}

type Space struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	DisplayName string `json:"displayName"`
}
