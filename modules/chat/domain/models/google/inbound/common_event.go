package models

type CommonEventObject struct {
	UserLocale string `json:"userLocale"`
	HostApp    string `json:"hostApp"`
	Platform   string `json:"platform"`
}