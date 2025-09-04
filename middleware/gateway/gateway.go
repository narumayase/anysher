package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	anysherhttp "github.com/narumayase/anysher/http"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
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
// - GATEWAY_IGNORE_ENDPOINTS -> format eg: GET:health|POST:send
// - LOG_LEVEL
func Sender() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := load()

		if !config.gatewayEnabled {
			c.Next()
			return
		}
		for _, endpoint := range config.ignoreEndpoints {
			// ignore the configurated endpoints
			if c.Request.Method == endpoint.Method && strings.Contains(c.Request.URL.Path, endpoint.Path) {
				c.Next()
				return
			}
		}
		bw := &bodyCaptureWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = bw

		c.Next()

		requestID := c.Request.Header.Get(requestIdHeader)
		if requestID == "" {
			// Generate a new one if not present
			requestID = uuid.NewString()
		}
		correlationID := c.Request.Header.Get(correlationIdHeader)
		routingID := c.Request.Header.Get(routingIdHeader)

		ctx := c.Request.Context()

		type Message struct {
			Content []byte `json:"content"`
		}
		payloadBytes, err := json.Marshal(Message{
			Content: bw.body.Bytes(),
		})
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to marshal response payload")
			return
		}

		httpClient := anysherhttp.NewClient(&http.Client{})
		resp, err := httpClient.Post(context.Background(), anysherhttp.Payload{
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
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			log.Ctx(ctx).Error().Err(
				fmt.Errorf("llm status code %d body %+v", resp.StatusCode, string(body))).
				Msg("failed to send response payload to gateway")
			return
		}
		log.Ctx(ctx).Info().Msg("response payload sent to gateway successfully")
	}
}

func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
