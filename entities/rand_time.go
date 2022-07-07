package entities

import (
	"fmt"
	"strings"
	"time"
)

type randTime time.Time

func (c *randTime) UnmarshalJSON(data []byte) error {
	const format = "2006-01-02 15:04:05Z"

	str := strings.Trim(strings.TrimSpace(string(data)), "\"")

	t, err := time.Parse(format, str)
	if err != nil {
		return fmt.Errorf("parse randTime: %w", err)
	}

	*c = randTime(t)

	return nil
}
