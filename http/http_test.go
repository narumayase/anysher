package http

import (
	"context"
	"encoding/json"
	config2 "github.com/narumayase/anysher/config"
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

		config := config2.Config{
			LogLevel: "info",
		}

		// Create a new HTTP client with the test server's URL
		client := NewClient(server.Client(), config)

		resp, err := client.Post(context.Background(), Payload{
			URL:     "",
			Token:   "",
			Headers: map[string]string{"Content-Type": "application/json"},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "success", respBody["message"])
	})

	t.Run("error during json.Marshal", func(t *testing.T) {
		config := config2.Config{
			LogLevel: "info",
		}
		client := NewClient(&http.Client{}, config)

		// Use a payload that cannot be marshaled to JSON (e.g., a channel)
		//	payload := make(chan int)
		//	headers := map[string]string{"Content-Type": "application/json"}

		resp, err := client.Post(context.Background(), Payload{
			URL:     "",
			Token:   "",
			Headers: map[string]string{"Content-Type": "application/json"},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal payload")
		assert.Nil(t, resp)
	})

	t.Run("error during http.NewRequest", func(t *testing.T) {
		config := config2.Config{
			LogLevel: "info",
		}
		client := NewClient(&http.Client{}, config)

		// Use an invalid URL to cause an error during NewRequest
		//payload := map[string]string{"test-key": "test-value"}
		//headers := map[string]string{"Content-Type": "application/json"}

		resp, err := client.Post(context.Background(), Payload{
			URL:     "",
			Token:   "",
			Headers: map[string]string{"Content-Type": "application/json"},
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
		mockClient := &http.Client{Transport: mockRT}
		client := NewHttpClient(mockClient, "test-token")

		payload := map[string]string{"test-key": "test-value"}
		headers := map[string]string{"Content-Type": "application/json"}

		resp, err := client.Post(context.Background(), headers, payload, "http://example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute request")
		assert.Nil(t, resp)
	})
}
