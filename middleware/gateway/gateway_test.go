package gateway_test

import (
	"bytes"
	"context"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/narumayase/anysher/middleware/gateway"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockClient implements el mismo método Post que anysherhttp.Client
type MockClient struct {
	Called   bool
	LastBody []byte
}

func (m *MockClient) Post(ctx context.Context, payload anysherhttp.Payload) (*http.Response, error) {
	m.Called = true
	m.LastBody = payload.Content
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil
}

func TestSenderMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := gateway.NewConfig(
		"debug",
		true,
		"http://mock",
		"token",
	)
	r := gin.New()
	r.Use(gateway.Sender(*config))
	r.POST("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	body := []byte(`{"foo":"bar"}`)
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
