package repository

import (
	"anysher/config"
	"anysher/internal/domain"
	"anysher/internal/infrastructure/client"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HTTPRepositoryImpl implements the domain.HTTPRepository interface using an HTTP client.
type HTTPRepositoryImpl struct {
	config     config.Config
	httpClient client.HttpClient
}

// NewHTTPRepository creates a new HTTPRepositoryImpl.
func NewHTTPRepository(
	config config.Config,
	httpClient client.HttpClient) domain.ProducerRepository {
	return &HTTPRepositoryImpl{
		httpClient: httpClient,
		config:     config,
	}
}

// Send sends the payload to the configured API endpoint.
func (h *HTTPRepositoryImpl) Send(ctx context.Context, payload domain.Payload) error {
	resp, err := h.httpClient.Post(ctx, payload.Headers, payload.Content, h.config.APIEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	log.Info().Msgf("API response status: %s", resp.Status)

	return nil
}

func (h *HTTPRepositoryImpl) Close() {
	log.Warn().Msg("HTTPRepositoryImpl closed")
}
