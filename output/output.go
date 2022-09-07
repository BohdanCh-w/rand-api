package output

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services"
)

const timeFormat = "15:04:05 02-01-2006"

var _ services.OutputGenerator = (*GeneratorImplementation)(nil)

func NewOutputProcessor(
	verbose bool,
	quiet bool,
	separator string,
	writer io.Writer,
) *GeneratorImplementation {
	return &GeneratorImplementation{
		verbose:   verbose,
		quiet:     quiet,
		separator: separator,
		writer:    writer,
	}
}

type GeneratorImplementation struct {
	verbose   bool
	quiet     bool
	separator string
	writer    io.Writer
}

func (svc *GeneratorImplementation) GenerateRandOutput(data []interface{}, apiInfo entities.APIInfo) error {
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

func (svc *GeneratorImplementation) generateAPIInfoOutput(apiInfo entities.APIInfo) {
	if svc.quiet {
		return
	}

	svc.showWarnings(apiInfo)

	if !svc.verbose {
		return
	}

	log.Printf("request %s finished at %s\n", apiInfo.ID.String(), apiInfo.Timestamp.Format(timeFormat))
	log.Printf("requests left: %d\n", apiInfo.RequestsLeft)
	log.Printf("random bits left: %d\n", apiInfo.BitsLeft)
	log.Printf("random bits used: %d\n", apiInfo.BitsUsed)
}

func (svc *GeneratorImplementation) showWarnings(apiInfo entities.APIInfo) {
	const (
		maxRequests = 1000
		maxBits     = 250_000

		percent      = 100
		warnTreshold = 5.0
	)

	requestsLeft := percent * float64(apiInfo.RequestsLeft) / maxRequests
	bitsLeft := percent * float64(apiInfo.BitsLeft) / maxBits

	if requestsLeft > warnTreshold && bitsLeft > warnTreshold {
		return
	}

	log.Printf("WARN: requests left    - %2.2f%% - %d\n", requestsLeft, apiInfo.RequestsLeft)
	log.Printf("WARN: random bits left - %2.2f%% - %d\n", requestsLeft, apiInfo.BitsLeft)
	log.Println()
}
