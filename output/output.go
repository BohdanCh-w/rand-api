package output

import (
	"fmt"
	"io"
	"os"
)

func GenerateOutput(cfg OutputConfig, data []interface{}) error {
	var w io.Writer = os.Stdout

	for _, v := range data {
		fmt.Fprintf(w, "%v%s", v, cfg.Separator)
	}

	return nil
}
