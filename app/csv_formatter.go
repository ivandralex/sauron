package sauron

import (
	"fmt"
	"io"
	"strings"
)

func printCSV(w io.Writer, row []string) error {
	sep := ""
	for _, cell := range row {
		_, err := fmt.Fprintf(w, `%s"%s"`, sep, strings.Replace(cell, `"`, `""`, -1))
		if err != nil {
			return err
		}
		sep = ","
	}
	_, err := fmt.Fprintf(w, "\n")
	if err != nil {
		return err
	}

	return nil
}
