package output

import (
	"fmt"
	"time"

	"github.com/bohdanch-w/rand-api/entities"
)

func (svc *OutputProcessorImplementation) GenerateUsageOutput(status entities.UsageStatus) error {
	format := `Usage statistic for API key %s:
  Status:        %s
  CreationTime:  %s
  TotalRequests: %d
  TotalBits:     %d
  RequestsLeft:  %d
  BitsLeft:      %d
`

	_, err := fmt.Fprintf(
		svc.writer,
		format,
		status.APIKey,
		status.Status,
		time.Time(status.CreationTime).Format(timeFormat),
		status.TotalRequests,
		status.TotalBits,
		status.RequestsLeft,
		status.BitsLeft,
	)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
