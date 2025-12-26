package dummy

import (
	"context"
	"encoding/json"
	"time"

	"logisync/internal/providers"
)

type Provider struct {
	name string
}

func New(name string) *Provider {
	return &Provider{name: name}
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Track(ctx context.Context, trackingCode string) (providers.Result, error) {
	_ = ctx
	payload := map[string]any{
		"provider":      p.Name(),
		"tracking_code": trackingCode,
		"status":        "IN_TRANSIT",
		"last_update":   time.Now().UTC().Format(time.RFC3339),
		"events": []map[string]any{
			{
				"timestamp":   time.Now().UTC().Format(time.RFC3339),
				"location":    "SAO PAULO - SP",
				"description": "Dummy tracking event",
			},
		},
		"raw": map[string]any{
			"source": "dummy",
		},
	}

	payloadBytes, _ := json.MarshalIndent(payload, "", "  ")
	artifacts := []providers.Artifact{
		{
			Kind:     "debug",
			Step:     "track",
			Filename: "payload.json",
			Data:     payloadBytes,
		},
	}

	return providers.Result{Payload: payload, Artifacts: artifacts}, nil
}
