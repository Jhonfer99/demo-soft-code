package infrastructure

import (
	rest "github.com/fsangopanta/demo-soft-code/modules/chat/infrastructure/google/rest"
	processor "github.com/fsangopanta/demo-soft-code/modules/chat/infrastructure/processor"
)

// =======================
// Controllers
// =======================

type GoogleController = rest.GoogleController

var NewGoogleController = rest.NewGoogleController

// =======================
// Processors
// =======================

type LocalProcessor = processor.LocalProcessor

func NewLocalProcessor() *LocalProcessor {
	return &processor.LocalProcessor{}
}
