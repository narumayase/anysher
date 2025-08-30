package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRoundTripper is a mock http.RoundTripper for testing purposes.
type MockRoundTripper struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

// RoundTrip implements the http.RoundTripper interface.
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// TestHttpClientImpl_Post tests the Post method of the Client.
func TestHttpClientImpl_Post(t *testing.T) {
	t.Run("successful POST request", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			var reqPayload map[string]string
			err := json.NewDecoder(r.Body).Decode(&reqPayload)
			assert.NoError(t, err)
			assert.Equal(t, "test-value", reqPayload["test-key"])
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"success"}`)) //nolint:errcheck
		}))
		defer server.Close()

		config := Config{
			LogLevel: "info",
		}

		// Create a new HTTP client with the test server's URL
		client := NewClient(server.Client(), config)

		// Create a sample payload
		payload := map[string]string{"test-key": "test-value"}
		payloadBytes, _ := json.Marshal(payload)

		resp, err := client.Post(context.Background(), Payload{
			URL:     server.URL,
			Token:   "test-token",
			Headers: map[string]string{"Content-Type": "application/json"},
			Content: payloadBytes,
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "success", respBody["message"])
	})

	

	t.Run("error during http.NewRequest", func(t *testing.T) {
		config := Config{
			LogLevel: "info",
		}
		client := NewClient(&http.Client{}, config)

		// Use an invalid URL to cause an error during NewRequest
		resp, err := client.Post(context.Background(), Payload{
			URL:     ":",
			Token:   "test-token",
			Headers: map[string]string{"Content-Type": "application/json"},
			Content: []byte("{}"),
		}) // Invalid URL
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create request")
		assert.Nil(t, resp)
	})

	t.Run("error during c.client.Do", func(t *testing.T) {
		// Create a mock RoundTripper that always returns an error
		mockRT := &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, assert.AnError
			},
		}

		config := Config{
			LogLevel: "info",
		}
		mockClient := &http.Client{Transport: mockRT}
		client := NewClient(mockClient, config)

		resp, err := client.Post(context.Background(), Payload{
			URL:     "http://example.com",
			Token:   "test-token",
			Headers: map[string]string{"Content-Type": "application/json"},
			Content: []byte("{}"),
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute request")
		assert.Nil(t, resp)
	})
}