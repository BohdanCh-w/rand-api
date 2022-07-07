package output

import (
	"fmt"
	"time"

	"github.com/bohdanch-w/rand-api/entities"
)

func GenerageStatusOutput(status entities.UsageStatus) {
	format := `Usage statistic for API key %s:
  Status:        %s
  CreationTime:  %s
  TotalRequests: %d
  TotalBits:     %d
  RequestsLeft:  %d
  BitsLeft:      %d
`

	fmt.Printf(
		format,
		status.APIKey,
		status.Status,
		time.Time(status.CreationTime).Format(timeFormat),
		status.TotalRequests,
		status.TotalBits,
		status.RequestsLeft,
		status.BitsLeft,
	)
}
