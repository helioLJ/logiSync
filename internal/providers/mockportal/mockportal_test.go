package mockportal

import (
	"context"
	"errors"
	"testing"

	"logisync/internal/providers"
)

func TestMapHTTPError(t *testing.T) {
	err := mapHTTPError(429, "rate")
	var providerErr *providers.Error
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected providers.Error")
	}
	if providerErr.Code != "RATE_LIMITED" {
		t.Fatalf("expected RATE_LIMITED, got %s", providerErr.Code)
	}

	err = mapHTTPError(404, "missing")
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected providers.Error")
	}
	if providerErr.Code != "INVALID_INPUT" {
		t.Fatalf("expected INVALID_INPUT, got %s", providerErr.Code)
	}
}

func TestTrackMissingBaseURL(t *testing.T) {
	provider := New(Config{BaseURL: ""})
	_, err := provider.Track(context.TODO(), "AA123")
	var providerErr *providers.Error
	if err == nil || !errors.As(err, &providerErr) {
		t.Fatalf("expected providers.Error")
	}
	if providerErr.Code != "INVALID_INPUT" {
		t.Fatalf("expected INVALID_INPUT, got %s", providerErr.Code)
	}
}

func TestAttachFailureArtifactsNilPage(t *testing.T) {
	provider := New(Config{BaseURL: "http://localhost"})
	err := provider.attachFailureArtifacts(nil, errors.New("boom"))
	var providerErr *providers.Error
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected providers.Error")
	}
	if providerErr.Code != "PROVIDER_ERROR" {
		t.Fatalf("expected PROVIDER_ERROR, got %s", providerErr.Code)
	}
	if len(providerErr.Artifacts) != 0 {
		t.Fatalf("expected no artifacts")
	}
}
