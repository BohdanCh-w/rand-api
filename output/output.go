package output

import (
	"fmt"
	"io"
	"os"

	"github.com/bohdanch-w/rand-api/entities"
)

func GenerateOutput(cfg OutputConfig, data []interface{}, apiInfo entities.APIInfo) error {
	if cfg.Verbose {
		generateAPIInfoOutput(apiInfo)
	}

	var w io.Writer = os.Stdout

	if cfg.Filename != "" {
		f, err := os.Create(cfg.Filename)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}

		w = f
	}

	for _, v := range data {
		fmt.Fprintf(w, "%v%s", v, cfg.Separator)
	}

	return nil
}

func generateAPIInfoOutput(apiInfo entities.APIInfo) {
	const timeFormat = "15:04:05 2006-01-02"

	fmt.Printf("request %s finished at %s\n", apiInfo.ID.String(), apiInfo.Timestamp.Format(timeFormat))
	fmt.Printf("requests -s")
}
