package main

import (
	"fmt"
	"os"

	"git.magenta.dk/os2datascanner/pdfanalyzer/pipeline"
)

// In case of error just print the error and panic.
func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// Main entry point for program.
// (Currently, it just serves as a small test program.)
func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Printf("ERROR - Invalid number of arguments: expected 1, got %d\n", len(args) - 1)
		os.Exit(-1)
	}
	
	file_name := args[1]
	pipe, err := pipeline.NewConcurrentPipeline(file_name)
	check(err)

	pipe.Run(
		pipeline.TextOnlyFilter,
		pipeline.SimpleExtractor,
		pipeline.CMapProcessor,
		pipeline.WritingReducer,
	)
}

