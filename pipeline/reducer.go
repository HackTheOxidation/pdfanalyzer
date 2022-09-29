package pipeline

import (
	"fmt"
	"os"

	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfobjects"
)

// Type signature for functions for the reducer step.
type ReducerFunction func (out []chan ProcessorResult, original *pdfobjects.Pdf)


// Collects all extracted data and prints it to STDIN.
func PrintingReducer(out []chan ProcessorResult, original *pdfobjects.Pdf) {
	strings := make([]string, 0)

	for i := range out {
		for obj := range out[i] {
			text := obj.stream
			if text != "" && obj.err == nil {
				strings = append(strings, text)
			}
		}
	}

	for _, str := range strings {
		if str != "" {
			fmt.Println(str)
		}
	}
}

// Collects all extracted data and writes it to a file.
func WritingReducer(out []chan ProcessorResult, original *pdfobjects.Pdf) {
	file_name := original.Name() + ".txt"
	fp, err := os.Create(file_name)
	check(err)

	defer fp.Close()

	for i := range out {
		for obj := range out[i] {
			text := obj.stream
			if text != "" && obj.err == nil {
				fp.WriteString(text)
			}
		}
	}
}
