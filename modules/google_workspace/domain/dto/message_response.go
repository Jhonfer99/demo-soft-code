package dto

type MessageResponse struct {
	Name string `json:"name,omitempty"`
	Text string `json:"text"`
	// User   UserResponse   `json:"user"`
	// Thread ThreadResponse `json:"thread"`
}

type UserResponse struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName"`
	// DomainId    string   `json:"domainId"`
	Type UserType `json:"type,omitempty"`
	// IsAnonymus  bool     `json:"isAnonymous"`
}

type ThreadResponse struct {
	Name string `json:"name,omitempty"`
}

type UserType string

const (
	TYPE_UNSPECIFIED UserType = "TYPE_UNSPECIFIED"
	HUMAN            UserType = "HUMAN"
	BOT              UserType = "BOT"
)