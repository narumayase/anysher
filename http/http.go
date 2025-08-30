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
func NewClient(client *http.Client, config Config) *Client {
	return &Client{
		client: client,
		config: config,
	}
}

// Post sends a POST request with JSON payload and bearer token authentication.
// It takes a context and a Payload struct as input.
// It returns the HTTP response and an error if the request fails.
func (c *Client) Post(ctx context.Context, payload Payload) (*http.Response, error) {
	payloadContent := payload.Content
	url := payload.URL
	headers := payload.Headers

	log.Debug().Msgf("payload to send: %s", string(payloadContent))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set custom headers from the payload.
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	log.Debug().Msgf("headers: to send to %s %+v", url, req.Header)

	// Set Authorization header with Bearer token.
	req.Header.Set("Authorization", "Bearer "+payload.Token)

	// Execute the HTTP request.
	resp, err := c.client.Do(req)
	if err != nil {
		log.Err(err).Msgf("Failed to send message via HTTP: %v", resp)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	log.Info().Msgf("API response %s status code: %d", url, resp.StatusCode)
	return resp, nil
}