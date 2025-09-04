package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	correlationIdHeader = "X-Correlation-ID"
	routingIdHeader     = "X-Routing-ID"
	requestIdHeader     = "X-Request-Id"
)

type bodyCaptureWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Sender middleware sends the request payload to the gateway after the handler has run.
// It takes the configuration from environment variables:
// - GATEWAY_API_URL
// - GATEWAY_ENABLED
// - GATEWAY_TOKEN
// - LOG_LEVEL
func Sender() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := load()

		if !config.gatewayEnabled {
			c.Next()
			return
		}
		ctx := c.Request.Context()
		bw := &bodyCaptureWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = bw

		c.Next()

		responseBody := bw.body.Bytes()

		httpClient := anysherhttp.NewClient(&http.Client{})

		requestID := c.Request.Header.Get(requestIdHeader)
		if requestID == "" {
			// Generate a new one if not present
			requestID = uuid.NewString()
		}
		correlationID := c.Request.Header.Get(correlationIdHeader)
		routingID := c.Request.Header.Get(routingIdHeader)

		type Message struct {
			Content []byte `json:"content"`
		}

		payloadBytes, err := json.Marshal(Message{
			Content: responseBody,
		})
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to marshal response payload")
			return
		}
		_, err = httpClient.Post(context.Background(), anysherhttp.Payload{
			URL:   config.gatewayAPIUrl,
			Token: config.gatewayToken,
			Headers: map[string]string{
				"Content-Type":      "application/json",
				correlationIdHeader: correlationID,
				routingIdHeader:     routingID,
				requestIdHeader:     requestID,
			},
			Content: payloadBytes,
		})
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to send response payload to gateway")
		} else {
			log.Ctx(ctx).Info().Msg("response payload sent to gateway successfully")
		}
	}
}

func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
