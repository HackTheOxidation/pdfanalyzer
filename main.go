package main

import (
	"fmt"
	"os"
)

// Main entry point for program.
// (Currently, it just serves as a small test program.)
func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Printf("ERROR - Invalid number of arguments: expected 1, got %d\n", len(args) - 1)
		os.Exit(-1)
	}
	
	file_name := args[1]
	reader, err := NewPdfReader(file_name)
	check(err)

	pdf := reader.ReadAll()

	fmt.Println("Name of file:", pdf.Name())
	fmt.Println("Number of objects read:", pdf.Count())

	obj, err := pdf.GetObject(0)
	check(err)

	fmt.Printf("First object: %v\n", obj)

	// content, err := obj.ExtractStream()
	// check(err)

	// fmt.Printf("Content of first stream: %s\n", content)

	fmt.Println("No errors reported.")
}

