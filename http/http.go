package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Payload represents the structure of an HTTP request payload.
type Payload struct {
	URL     string
	Token   string
	Headers map[string]string
	Content []byte
}

// Client represents an HTTP client with a base client and configuration.
type Client struct {
	client *http.Client
	config Config
}

// NewClient creates a new HTTP client with bearer token authentication.
// It takes an http.Client and a Config struct as input.
// It takes the environment variable LOG_LEVEL
func NewClient(client *http.Client) *Client {
	// load configuration from environment
	cfg := load()

	return &Client{
		client: client,
		config: cfg,
	}
}

// Post sends a POST request with JSON payload and bearer token authentication.
// It takes a context and a Payload struct as input.
// It returns the HTTP response and an error if the request fails.
func (c *Client) Post(ctx context.Context, payload Payload) (*http.Response, error) {
	payloadContent := payload.Content
	url := payload.URL
	headers := payload.Headers

	log.Ctx(ctx).Debug().Msgf("payload to send to %s: %s", url, string(payloadContent))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set custom headers from the payload.
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	log.Ctx(ctx).Debug().Msgf("headers: to send to %s %+v", url, req.Header)

	// Set Authorization header with Bearer token.
	req.Header.Set("Authorization", "Bearer "+payload.Token)

	// Execute the HTTP request.
	resp, err := c.client.Do(req)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("Failed to send message to %s via HTTP: %v", url, resp)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	log.Ctx(ctx).Info().Msgf("API response %s status code: %d", url, resp.StatusCode)
	return resp, nil
}
