package providers

import (
	"context"
)

type Artifact struct {
	Kind     string
	Step     string
	Filename string
	Data     []byte
}

type Result struct {
	Payload   map[string]any
	Artifacts []Artifact
}

type Provider interface {
	Name() string
	Track(ctx context.Context, trackingCode string) (Result, error)
}

type Error struct {
	Code      string
	Message   string
	Err       error
	Artifacts []Artifact
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Message
	}
	if e.Err == e {
		return e.Message
	}
	return e.Message + ": " + e.Err.Error()
}
