package domain

import (
	"context"
)

// ProducerRepository defines the interface for the producer repository
type ProducerRepository interface {
	Send(ctx context.Context, payload Payload) error
	Close()
}
