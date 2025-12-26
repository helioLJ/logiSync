package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DB_URL", "postgres://test")
	t.Setenv("OP_TIMEOUT", "5s")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected default HTTP addr, got %s", cfg.HTTPAddr)
	}
	if cfg.OpTimeout != 5*time.Second {
		t.Fatalf("expected default timeout, got %s", cfg.OpTimeout)
	}
}

func TestEnvBool(t *testing.T) {
	t.Setenv("PLAYWRIGHT_HEADLESS", "false")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.PlaywrightHeadless {
		t.Fatalf("expected headless false")
	}
}

func TestEnvBoolInvalidFallback(t *testing.T) {
	t.Setenv("DB_URL", "postgres://test")
	t.Setenv("PLAYWRIGHT_HEADLESS", "maybe")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !cfg.PlaywrightHeadless {
		t.Fatalf("expected headless true fallback")
	}
}

func TestEnvDurationInvalidFallback(t *testing.T) {
	t.Setenv("DB_URL", "postgres://test")
	t.Setenv("OP_TIMEOUT", "not-a-duration")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.OpTimeout != 5*time.Second {
		t.Fatalf("expected fallback timeout, got %s", cfg.OpTimeout)
	}
}
