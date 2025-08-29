package domain

type Payload struct {
	KafkaPayload
	HTTPPayload

	Headers map[string]string
	Content []byte
}

type HTTPPayload struct {
	URL   string
	Token string
}

type KafkaPayload struct {
	Key string
}

// ErrorResponse represents an error model
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
