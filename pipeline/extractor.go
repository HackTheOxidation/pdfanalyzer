package pipeline

import (
	"fmt"
	"sync"

	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfobjects"
)

// Type signature for functions for the extraction step.
type ExtractorFunction func (in <-chan pdfobjects.PdfObject, out chan<- ExtractorResult, cmap *pdfobjects.CMap)

// Wrapper for running the extractor stage with the extractor function `f`.
func runExtractorStage(
	f ExtractorFunction,
	in <-chan pdfobjects.PdfObject,
	out chan<- ExtractorResult,
	cmap *pdfobjects.CMap,
	wg *sync.WaitGroup,
) {
	f(in, out, cmap)
	wg.Done()
}

// Simple extractor function that extracts text and cmaps.
func SimpleExtractor(in <-chan pdfobjects.PdfObject, out chan<- ExtractorResult, cmap *pdfobjects.CMap) {
	defer close(out)
	for data := range in {
		out <- NewExtractorResult(data.ExtractStream(cmap))
	}
}

type ExtractorResult struct {
	stream string
	process bool
	err error
}

func NewExtractorResult(stream string, process bool, err error) ExtractorResult {
	return ExtractorResult{
		stream,
		process,
		err,
	}
}

func (e ExtractorResult) ToProcessorResult() ProcessorResult {
	return NewProcessorResult(e.stream, e.err)
}

func (e ExtractorResult) String() string {
	return fmt.Sprintf("stream: %s", e.stream)
}
