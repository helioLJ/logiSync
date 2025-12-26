package mockportal

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"

	"logisync/internal/providers"
)

type Config struct {
	BaseURL  string
	Timeout  time.Duration
	Headless bool
	SlowMo   time.Duration // Slow motion delay between actions (for debugging)
}

type Provider struct {
	cfg Config
}

type trackResponse struct {
	TrackingCode string  `json:"tracking_code"`
	Status       string  `json:"status"`
	LastUpdate   string  `json:"last_update"`
	Events       []event `json:"events"`
}

type event struct {
	Timestamp   string `json:"timestamp"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

func New(cfg Config) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Name() string {
	return "mock_portal_scrape"
}

func (p *Provider) Track(ctx context.Context, trackingCode string) (providers.Result, error) {
	if strings.TrimSpace(p.cfg.BaseURL) == "" {
		return providers.Result{}, &providers.Error{Code: "INVALID_INPUT", Message: "missing mock portal url"}
	}

	pw, err := playwright.Run()
	if err != nil {
		return providers.Result{}, &providers.Error{Code: "PROVIDER_ERROR", Message: "failed to start playwright", Err: err}
	}
	defer pw.Stop()

	launchOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(p.cfg.Headless),
	}
	if p.cfg.SlowMo > 0 {
		launchOpts.SlowMo = playwright.Float(float64(p.cfg.SlowMo.Milliseconds()))
	}
	browser, err := pw.Chromium.Launch(launchOpts)
	if err != nil {
		return providers.Result{}, &providers.Error{Code: "PROVIDER_ERROR", Message: "failed to launch browser", Err: err}
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return providers.Result{}, &providers.Error{Code: "PROVIDER_ERROR", Message: "failed to open page", Err: err}
	}

	result, err := p.runFlow(ctx, page, trackingCode)
	if err != nil {
		if providerErr := p.attachFailureArtifacts(page, err); providerErr != nil {
			return providers.Result{}, providerErr
		}
		return providers.Result{}, err
	}

	return result, nil
}

func (p *Provider) runFlow(ctx context.Context, page playwright.Page, trackingCode string) (providers.Result, error) {
	timeoutMs := float64(p.cfg.Timeout.Milliseconds())
	if timeoutMs <= 0 {
		timeoutMs = 30000 // Default 30s para Playwright
	}

	url := strings.TrimRight(p.cfg.BaseURL, "/") + "/track"
	if _, err := page.Goto(url, playwright.PageGotoOptions{Timeout: playwright.Float(timeoutMs)}); err != nil {
		return providers.Result{}, &providers.Error{Code: "PROVIDER_ERROR", Message: "failed to load portal", Err: err}
	}

	if err := page.Fill("[data-testid=track-input]", trackingCode, playwright.PageFillOptions{Timeout: playwright.Float(timeoutMs)}); err != nil {
		return providers.Result{}, &providers.Error{Code: "PROVIDER_ERROR", Message: "failed to fill tracking code", Err: err}
	}

	// Usar padrão glob do Playwright para capturar a resposta da API
	// O padrão **/api/track/* funciona para qualquer URL que contenha /api/track/
	apiURLPattern := "**/api/track/*"
	resp, err := page.ExpectResponse(apiURLPattern, func() error {
		return page.Click("[data-testid=track-submit]", playwright.PageClickOptions{Timeout: playwright.Float(timeoutMs)})
	}, playwright.PageExpectResponseOptions{Timeout: playwright.Float(timeoutMs)})
	if err != nil {
		return providers.Result{}, &providers.Error{Code: "TIMEOUT", Message: "timed out waiting for response", Err: err}
	}

	body, err := resp.Body()
	if err != nil {
		return providers.Result{}, &providers.Error{Code: "PROVIDER_ERROR", Message: "failed to read response", Err: err}
	}

	if resp.Status() >= 400 {
		return providers.Result{}, mapHTTPError(int(resp.Status()), string(body))
	}

	var apiResp trackResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return providers.Result{}, &providers.Error{Code: "PARSE_ERROR", Message: "failed to parse response", Err: err}
	}

	events := make([]map[string]any, 0, len(apiResp.Events))
	for _, evt := range apiResp.Events {
		events = append(events, map[string]any{
			"timestamp":   evt.Timestamp,
			"location":    evt.Location,
			"description": evt.Description,
		})
	}

	payload := map[string]any{
		"provider":      p.Name(),
		"tracking_code": apiResp.TrackingCode,
		"status":        apiResp.Status,
		"last_update":   apiResp.LastUpdate,
		"events":        events,
		"raw": map[string]any{
			"response": json.RawMessage(body),
		},
	}

	return providers.Result{
		Payload: payload,
		Artifacts: []providers.Artifact{
			{
				Kind:     "response",
				Step:     "track",
				Filename: "response.json",
				Data:     body,
			},
		},
	}, nil
}

func (p *Provider) attachFailureArtifacts(page playwright.Page, err error) error {
	providerErr := &providers.Error{Code: "PROVIDER_ERROR", Message: "tracking failed", Err: err}
	var existing *providers.Error
	if errors.As(err, &existing) {
		providerErr = existing
		if providerErr.Err == nil && providerErr != err {
			providerErr.Err = err
		}
	}
	if page == nil {
		return providerErr
	}

	artifacts := []providers.Artifact{}
	screenshot, shotErr := page.Screenshot(playwright.PageScreenshotOptions{FullPage: playwright.Bool(true)})
	if shotErr == nil {
		artifacts = append(artifacts, providers.Artifact{
			Kind:     "screenshot",
			Step:     "track",
			Filename: "failure.png",
			Data:     screenshot,
		})
	}

	html, htmlErr := page.Content()
	if htmlErr == nil {
		artifacts = append(artifacts, providers.Artifact{
			Kind:     "html",
			Step:     "track",
			Filename: "failure.html",
			Data:     []byte(html),
		})
	}

	if len(artifacts) > 0 {
		providerErr.Artifacts = artifacts
	}

	return providerErr
}

func mapHTTPError(status int, message string) error {
	code := "PROVIDER_ERROR"
	switch status {
	case 400, 404:
		code = "INVALID_INPUT"
	case 401, 403:
		code = "AUTH_ERROR"
	case 408, 504:
		code = "TIMEOUT"
	case 429:
		code = "RATE_LIMITED"
	}
	return &providers.Error{Code: code, Message: strings.TrimSpace(message)}
}
