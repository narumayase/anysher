package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HTTPClientImpl HttpClientImpl implements HttpClient interface for making HTTP requests
type HTTPClientImpl struct {
	client      *http.Client
	bearerToken string
}

// HttpClient defines the interface for making HTTP requests
type HttpClient interface {
	Post(ctx context.Context, headers map[string]string, payload interface{}, url string) (*http.Response, error)
}

// NewHttpClient creates a new HTTP client with bearer token authentication
func NewHttpClient(client *http.Client, bearerToken string) HttpClient {
	return &HTTPClientImpl{
		client:      client,
		bearerToken: bearerToken,
	}
}

// Post sends a POST request with JSON payload and bearer token authentication
func (c *HTTPClientImpl) Post(ctx context.Context, headers map[string]string, payload interface{}, url string) (*http.Response, error) {
	// Convert payload to JSON
	var jsonPayload []byte
	switch v := payload.(type) {
	case []byte:
		jsonPayload = v
	default:
		var err error
		jsonPayload, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}
	log.Debug().Msgf("payload to send: %s", string(jsonPayload))
	log.Debug().Msgf("url %s", url)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		// Set headers
		req.Header.Set(key, value)
	}
	log.Debug().Msgf("headers: to send to %s %+v", url, req.Header)

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return resp, nil
}
