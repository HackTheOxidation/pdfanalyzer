package main

import (
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfanalyzer"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	file_name := "assets/gcc.pdf"
	reader, err := pdfanalyzer.NewPdfReader(file_name)
	check(err)

	reader.ReadAll()
}

