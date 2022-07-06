package output

import (
	"fmt"
	"io"
	"os"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
)

func GenerateOutput(cfg config.Output, data []interface{}, apiInfo entities.APIInfo) error {
	generateAPIInfoOutput(cfg, apiInfo)

	var w io.Writer = os.Stdout

	if cfg.Filename != nil {
		f, err := os.Create(*cfg.Filename)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}

		w = f
	}

	for i, v := range data {
		fmt.Fprintf(w, "%v", v)

		if i+1 != len(data) {
			fmt.Fprint(w, cfg.Separator)
		}
	}

	return nil
}

func generateAPIInfoOutput(cfg config.Output, apiInfo entities.APIInfo) {
	if cfg.Quite {
		return
	}

	showWarnings(apiInfo)

	if !cfg.Verbose {
		return
	}

	const timeFormat = "15:04:05 2006-01-02"

	fmt.Printf("request %s finished at %s\n", apiInfo.ID.String(), apiInfo.Timestamp.Format(timeFormat))
	fmt.Printf("requests left: %d\n", apiInfo.RequestsLeft)
	fmt.Printf("random bits left: %d\n", apiInfo.BitsLeft)
	fmt.Printf("random bits used: %d\n", apiInfo.BitsUsed)
}

func showWarnings(apiInfo entities.APIInfo) {
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
