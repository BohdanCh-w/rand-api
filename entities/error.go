package entities

import "github.com/google/uuid"

type ErrorResponse struct {
	ID             uuid.UUID `json:"id"`
	JsonrpcVersion string    `json:"jsonrpc"`
	Error          struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
