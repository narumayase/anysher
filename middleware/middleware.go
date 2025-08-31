package middleware

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const requestIdHeader = "X-Request-Id"

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	defaultConfig := cors.DefaultConfig()
	defaultConfig.AllowAllOrigins = true
	defaultConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	defaultConfig.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Routing-Id",
		"X-Correlation-Id",
		"X-Request-Id",
	}
	defaultConfig.ExposeHeaders = []string{"Content-Length"}
	defaultConfig.AllowCredentials = true
	defaultConfig.MaxAge = 12 * time.Hour

	return cors.New(defaultConfig)
}

// Logger middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		ctx := param.Request.Context()

		log.Ctx(ctx).Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Dur("latency", param.Latency).
			Str("client_ip", param.ClientIP).
			Str("user_agent", param.Request.UserAgent()).
			Msg("HTTP Request")
		return ""
	})
}

// ErrorHandler middleware for handling panics and errors
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(500, ErrorResponse{
				Error:   "internal_server_error",
				Message: err,
				Code:    500,
			})
		} else {
			c.JSON(500, ErrorResponse{
				Error:   "internal_server_error",
				Message: "An unexpected error occurred",
				Code:    500,
			})
		}
	})
}

// HeadersToContext is a Gin middleware that takes every incoming request header
// and stores it individually into the request context.
// This allows you to access any header later in the request lifecycle using ctx.Value(headerName).
func HeadersToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Iterate through all headers and inject them into the context
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				// Store each header's first value under its name
				ctx = context.WithValue(ctx, k, v[0])
			}
		}
		log.Ctx(ctx).Info().Msgf("headers received: %+v", c.Request.Header)
		// Replace the request with the new one that has the updated context
		c.Request = c.Request.WithContext(ctx)

		// Continue to the next middleware/handler
		c.Next()
	}
}

// RequestIDToLogger is a middleware that extracts the X-Request-Id header
// (or generates a new UUID if missing) and injects it into zerolog's context.
// Any log written with log.Ctx(ctx) will automatically include "request_id".
func RequestIDToLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Look for X-Request-Id header
		requestID := c.Request.Header.Get(requestIdHeader)
		if requestID == "" {
			// Generate a new one if not present
			requestID = uuid.NewString()
		}

		// Create a logger with request_id and attach it to the context
		logger := log.With().Str("request_id", requestID).Logger()
		ctx = logger.WithContext(ctx)

		// Replace the request with the new context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
