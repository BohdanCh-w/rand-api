package entities

import (
	"time"

	"github.com/google/uuid"
)

type APIInfo struct {
	ID           uuid.UUID
	Timestamp    time.Time
	BitsUsed     uint64
	BitsLeft     uint64
	RequestsLeft uint64
}
