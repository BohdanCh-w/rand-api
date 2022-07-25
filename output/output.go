package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services"
)

const timeFormat = "15:04:05 02-01-2006"

var _ services.OutputProcessor = (*OutputProcessorImplementation)(nil)

func NewOutputProcessor(
	verbose bool,
	quiet bool,
	separator string,
	writer io.Writer,
) *OutputProcessorImplementation {
	return &OutputProcessorImplementation{
		verbose:   verbose,
		quiet:     quiet,
		separator: separator,
		writer:    writer,
	}
}

type OutputProcessorImplementation struct {
	verbose   bool
	quiet     bool
	separator string
	writer    io.Writer
}

func (svc *OutputProcessorImplementation) GenerateRandOutput(data []interface{}, apiInfo entities.APIInfo) error {
	svc.generateAPIInfoOutput(apiInfo)

	dataStr := make([]string, 0, len(data))

	for _, v := range data {
		dataStr = append(dataStr, fmt.Sprintf("%v", v))
	}

	if _, err := fmt.Fprintln(svc.writer, strings.Join(dataStr, svc.separator)); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func (svc *OutputProcessorImplementation) generateAPIInfoOutput(apiInfo entities.APIInfo) {
	if svc.quiet {
		return
	}

	svc.showWarnings(apiInfo)

	if !svc.verbose {
		return
	}

	fmt.Printf("request %s finished at %s\n", apiInfo.ID.String(), apiInfo.Timestamp.Format(timeFormat))
	fmt.Printf("requests left: %d\n", apiInfo.RequestsLeft)
	fmt.Printf("random bits left: %d\n", apiInfo.BitsLeft)
	fmt.Printf("random bits used: %d\n", apiInfo.BitsUsed)
}

func (svc *OutputProcessorImplementation) showWarnings(apiInfo entities.APIInfo) {
	const (
		maxRequests = 1000
		maxBits     = 250_000
	)

	requestsLeft := 100 * float64(apiInfo.RequestsLeft) / maxRequests
	bitsLeft := 100 * float64(apiInfo.BitsLeft) / maxBits

	if requestsLeft > 5.0 && bitsLeft > 5.0 {
		return
	}

	fmt.Printf("WARN: requests left    - %2.2f%% - %d\n", requestsLeft, apiInfo.RequestsLeft)
	fmt.Printf("WARN: random bits left - %2.2f%% - %d\n", requestsLeft, apiInfo.BitsLeft)
	fmt.Println()
}
