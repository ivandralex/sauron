package writers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//CSVWriter writes features and label to csv file
type CSVWriter struct {
	file *os.File
}

//Init initializes CSVWriter by opening file
func (w *CSVWriter) Init(path string) {
	absPath, _ := filepath.Abs(path)
	os.Remove(absPath)

	file, err := os.Create(absPath)

	if err != nil {
		log.Fatalln("Error creating features file:", err)
	}

	if err != nil {
		log.Fatalln("Could not open dump file for writing:", err)
	}

	w.file = file
}

//WriteHead writes csv file header
func (w *CSVWriter) WriteHead(featureNames []string) {
	columnNames := []string{"key"}
	columnNames = append(columnNames, featureNames...)
	columnNames = append(columnNames, "label")

	if err := printCSV(w.file, columnNames); err != nil {
		log.Fatalln("Error writing header to csv:", err)
	}
}

//WriteSession writes session
func (w *CSVWriter) WriteSession(key string, features []string, label string) {
	line := []string{key}
	line = append(line, features...)
	line = append(line, label)

	if err := printCSV(w.file, line); err != nil {
		log.Fatalln("Error writing session to csv:", err)
	}
}

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
