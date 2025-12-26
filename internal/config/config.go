package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr           string
	DBURL              string
	RedisAddr          string
	RedisStream        string
	RedisGroup         string
	RedisConsumer      string
	OpTimeout          time.Duration
	ArtifactsRoot      string
	MockPortalURL      string
	PlaywrightHeadless bool
	PlaywrightSlowMo   time.Duration
}

func Load() (Config, error) {
	// Load dotenv file if it exists
	_ = godotenv.Load()

	cfg := Config{
		HTTPAddr:           env("HTTP_ADDR", ":8080"),
		DBURL:              env("DB_URL", "postgres://postgres:postgres@localhost:5432/logisync?sslmode=disable"),
		RedisAddr:          env("REDIS_ADDR", "localhost:6379"),
		RedisStream:        env("REDIS_STREAM", "tracking:jobs"),
		RedisGroup:         env("REDIS_GROUP", "tracking-workers"),
		RedisConsumer:      env("REDIS_CONSUMER", "worker-1"),
		OpTimeout:          envDuration("OP_TIMEOUT", 5*time.Second),
		ArtifactsRoot:      env("ARTIFACTS_ROOT", "./artifacts"),
		MockPortalURL:      env("MOCK_PORTAL_URL", "http://localhost:8090"),
		PlaywrightHeadless: envBool("PLAYWRIGHT_HEADLESS", true),
		PlaywrightSlowMo:   envDuration("PLAYWRIGHT_SLOW_MO", 0),
	}

	if cfg.DBURL == "" {
		return Config{}, fmt.Errorf("DB_URL is required")
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if val := strings.ToLower(strings.TrimSpace(os.Getenv(key))); val != "" {
		switch val {
		case "1", "true", "yes", "y":
			return true
		case "0", "false", "no", "n":
			return false
		}
	}
	return fallback
}
