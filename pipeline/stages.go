package pipeline

import (
	"fmt"
	"runtime"
	"sync"

	"git.magenta.dk/os2datascanner/pdfanalyzer/parser"
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfobjects"
)

// General interface for a pipeline.
type Pipeline interface {
	 Run(FilterFunction, ExtractorFunction, ProcessorFunction, ReducerFunction)
}

// Container for pipeline.
type ConcurrentPipeline struct {
	reader *parser.PdfReader
	pdf *pdfobjects.Pdf
}

type indexRange struct {
	First int
	Last int
}

// Returns a new Pipeline struct.
func NewConcurrentPipeline(file_name string) (Pipeline, error) {
	reader, err := parser.NewPdfReader(file_name)
	return ConcurrentPipeline{
		reader,
		pdfobjects.NewPdf(file_name),
	}, err
}

// Create index ranges to partition the data with.
func (p ConcurrentPipeline) partition(cores int) []indexRange {
	ranges := make([]indexRange, cores)

	step := p.pdf.Count() / cores

	for i := 0; i < cores; i++ {
		first := i * step
		last := (i + 1) * step - 1
		ranges[i] = indexRange{ First: first, Last: last }
	}

	return ranges
}

func generateChannels(cores int, nobjects int)(
	[]chan pdfobjects.PdfObject,
	[]chan pdfobjects.PdfObject,
	[]chan ExtractorResult,
	[]chan ProcessorResult,
	[]chan bool,
) {
	gen := make([]chan pdfobjects.PdfObject, cores)
	fil := make([]chan pdfobjects.PdfObject, cores)
	ext := make([]chan ExtractorResult, cores)
	pro := make([]chan ProcessorResult, cores)
	fin := make([]chan bool, cores)

	bufferSize := nobjects / cores + nobjects % cores

	for i := 0; i < cores; i++ {
		gen[i] = make(chan pdfobjects.PdfObject)
		fil[i] = make(chan pdfobjects.PdfObject)
		ext[i] = make(chan ExtractorResult, bufferSize)
		pro[i] = make(chan ProcessorResult)
		fin[i] = make(chan bool)
	}

	return gen, fil, ext, pro, fin
}

func fillChannel(pdf *pdfobjects.Pdf, out chan pdfobjects.PdfObject, index indexRange) {
	defer close(out)

	for i := index.First; i < index.Last; i++ {
		data, err := pdf.GetObject(i)
		check(err)

		out <- *data
	}
}

// Run a concurrent pipeline.
func (p ConcurrentPipeline) Run(
	filter FilterFunction,
	extract ExtractorFunction,
	process ProcessorFunction,
	reduce ReducerFunction,
) {
	fmt.Println("Running Pipeline for file:", p.pdf.Name())

	// Read and parse the pdf document using a single core.
	p.pdf = p.reader.ReadAll()

	// Get the number of available cores on the system.
	cores := runtime.NumCPU()

	// Partition the objects into index ranges and generate channels 
	ranges := p.partition(cores)
	gen, fil, ext, pro, _ := generateChannels(cores, p.pdf.Count())

	// Make a global CMap
	cmap := make(pdfobjects.CMap, 0)

	// Initialize a waitgroup
	var wg sync.WaitGroup
	wg.Add(cores)

	// Start filling channels, filtering and extraction as goroutines.
	for i := 0; i < cores; i++ {
		go fillChannel(p.pdf, gen[i], ranges[i])
		go runFilterStage(filter, gen[i], fil[i])
		go runExtractorStage(extract, fil[i], ext[i], &cmap, &wg)
	}

	// Wait for extraction to finish in case cmaps are last.
	wg.Wait()
	
	// Start processing stage as goroutines
	for i := 0; i < cores; i++ {
		go runProcessorStage(process, ext[i], pro[i], &cmap)
	}

	// Reduce the result.
	reduce(pro, p.pdf)

	fmt.Println("Pipeline ran successfully! No errors reported.")
}

// In case of error just print the error and panic.
func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
