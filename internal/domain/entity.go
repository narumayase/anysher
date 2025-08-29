package domain

type Payload struct {
	Key     string
	Headers map[string]string
	Content []byte
}

// ErrorResponse represents an error model
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
