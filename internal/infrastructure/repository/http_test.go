package repository

import (
	"anysher/config"
	"anysher/internal/domain"
	clientmocks "anysher/internal/infrastructure/client/mocks"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendRepositoryImpl_Send(t *testing.T) {
	mockHTTPClient := new(clientmocks.MockHTTPClient)
	cfg := config.Config{APIEndpoint: "http://localhost:8080"}
	repo := NewHTTPRepository(cfg, mockHTTPClient)

	ctx := context.Background()
	msg := domain.Payload{Content: []byte(`{"key":"value"}`)}

	t.Run("success", func(t *testing.T) {
		mockResponse := clientmocks.CreateMockResponse(http.StatusOK, `{"status":"ok"}`)
		mockHTTPClient.On("Post", ctx, mock.AnythingOfType("map[string]string"), msg.Content, cfg.APIEndpoint).Return(mockResponse, nil).Once()

		err := repo.Send(ctx, msg)

		assert.NoError(t, err)
		mockHTTPClient.AssertExpectations(t)
	})

	t.Run("http client error", func(t *testing.T) {
		expectedErr := errors.New("http client error")
		mockHTTPClient.On("Post", ctx, mock.AnythingOfType("map[string]string"), msg.Content, cfg.APIEndpoint).Return(nil, expectedErr).Once()

		err := repo.Send(ctx, msg)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockHTTPClient.AssertExpectations(t)
	})

	t.Run("unexpected status code", func(t *testing.T) {
		mockResponse := clientmocks.CreateMockResponse(http.StatusInternalServerError, `{"error":"internal server error"}`)
		mockHTTPClient.On("Post", ctx, mock.AnythingOfType("map[string]string"), msg.Content, cfg.APIEndpoint).Return(mockResponse, nil).Once()

		err := repo.Send(ctx, msg)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
		mockHTTPClient.AssertExpectations(t)
	})

	t.Run("success with correlation ID header", func(t *testing.T) {
		correlationID := "test-correlation-id"
		msgWithHeader := domain.Payload{
			Content: []byte(`{"key":"value"}`),
			Headers: map[string]string{"correlation_id": correlationID},
		}
		mockResponse := clientmocks.CreateMockResponse(http.StatusOK, `{"status":"ok"}`)
		mockHTTPClient.On("Post", ctx, mock.MatchedBy(func(h map[string]string) bool { return h["X-Correlation-ID"] == correlationID }), msgWithHeader.Content, cfg.APIEndpoint).Return(mockResponse, nil).Once()

		err := repo.Send(ctx, msgWithHeader)

		assert.NoError(t, err)
		mockHTTPClient.AssertExpectations(t)
	})
}
