package models

type WorkspaceEvent struct {
	Common CommonEventObject `json:"commonEventObject"`
	Chat   *ChatEvent        `json:"chat"`
}