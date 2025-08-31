package gateway

import (
	"bytes"
	"github.com/gin-gonic/gin"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

const (
	correlationIdHeader = "X-Correlation-ID"
	routingIdHeader     = "X-Routing-ID"
	requestIdHeader     = "X-Request-Id"
)

// Sender middleware sends the request payload to the gateway after the handler has run.
func Sender(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If gateway is disabled, skip sending
		if !config.gatewayEnabled {
			c.Next()
			return
		}
		cfg := anysherhttp.NewConfiguration(config.logLevel)
		ctx := c.Request.Context()

		// Create a new HTTP client
		httpClient := anysherhttp.NewClient(&http.Client{}, cfg)

		// Read request body (store it for later)
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to read request body for gateway")
			c.Next()
			return
		}

		// Restore body so the handler can still read it
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// Run handler first
		c.Next()

		// After handler finished, send to gateway
		correlationID, _ := ctx.Value(correlationIdHeader).(string)
		routingID, _ := ctx.Value(routingIdHeader).(string)
		requestID, _ := ctx.Value(requestIdHeader).(string)

		_, postErr := httpClient.Post(ctx, anysherhttp.Payload{
			URL:   config.gatewayAPIUrl,
			Token: config.gatewayToken,
			Headers: map[string]string{
				"Content-Type":      "application/json",
				correlationIdHeader: correlationID,
				routingIdHeader:     routingID,
				requestIdHeader:     requestID,
			},
			Content: body,
		})
		if postErr != nil {
			log.Ctx(ctx).Error().Err(postErr).Msg("failed to send payload to gateway")
		} else {
			log.Ctx(ctx).Info().Msg("payload sent to gateway successfully")
		}
	}
}
