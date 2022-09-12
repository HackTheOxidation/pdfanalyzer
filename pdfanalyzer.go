package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfobjects"
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdftypes"
)

// A Reader for pdf files.
//
// Struct for reading pdf files by using and internal `refreshingReader`
type PdfReader struct {
	filename string
	reader *refreshingReader
}

// Constructs a new PdfReader with `filename` and a `refreshingReader`.
func NewPdfReader(filename string) (*PdfReader, error) {
	reader, err := newRefreshingReader(filename)
	return &PdfReader{
		filename,
		reader,
	}, err
}

// In case of error, close `PdfReader` and terminate program.
func continueOrClose(err error, r *PdfReader) {
	if err != nil {
		fmt.Println(err)
		r.close()
		os.Exit(-1)
	}
}

// Reads an entire pdf file and return a `Pdf` struct.
func (r *PdfReader) ReadAll() *pdfobjects.Pdf {
	pdf := pdfobjects.NewPdf(r.filename)
	line_number := 0

	// Try parsing objects as long as End-Of-File (EOF) is not reached.
	for !r.reader.IsEOF() {
		line, _, err := r.reader.ReadLine()
		line_number++
		continueOrClose(err, r)
		line_str := string(line)

		if pdftypes.ObjectBegins(line_str) {
			obj := r.readObject(line_number)
			pdf.AppendObject(obj)
		} else if pdftypes.IsVersion(line_str) {
			// Parse PDF version.
			pdf.SetVersion(line_str)
		}
	}
	
	r.close()

	return pdf
}

// Close the readers file handle.
func (r *PdfReader) close() {
	r.reader.close()
}

// Internal buffering reader. Reads 4kB at a time to keep memory low.
// The buffer can be reloaded.
type refreshingReader struct {
	fp *os.File
	reader *bufio.Reader
}

// Construct a new refreshingReader with `filename`.
func newRefreshingReader(filename string) (*refreshingReader, error) {
	fp, err := os.Open(filename)
	reader := bufio.NewReader(fp)
	return &refreshingReader{
		fp,
		reader,
	}, err
}

// Read the next 4kB from file into buffer.
// Current buffer is discarded.
func (r *refreshingReader) reloadBuffer() {
	r.reader = bufio.NewReader(r.fp)
}

// Close the file handle.
func (r *refreshingReader) close() {
	r.fp.Close()
}

// Get the name of the file.
func (r *refreshingReader) name() string {
	return r.fp.Name()
}

// Reads one byte from the file.
func (r *refreshingReader) ReadByte() (byte, error) {
	b, err := r.reader.ReadByte()
	if err == io.EOF {
		r.reloadBuffer()
		return r.ReadByte()
	}
	return b, err
}

// Read `n` bytes without moving the cursor of the reader.
func (r *refreshingReader) Peek(n int) ([]byte, error) {
	return r.reader.Peek(n)
}

// Reads a line from the file.
func (r *refreshingReader) ReadLine() ([]byte, bool, error) {
	return r.reader.ReadLine()
}

// Checks whether the reader has reached End-Of-File (EOF).
func (r *refreshingReader) IsEOF() bool {
	peek, err := r.reader.Peek(pdftypes.EOFL)
	check(err)
	return pdftypes.IsEOF(string(peek))
}

// Check whether the reader has reached the end of a stream.
func (r *refreshingReader) IsEndstream() bool {
	peek, err := r.Peek(pdftypes.ENDSTREAML)
	check(err)
	return pdftypes.StreamEnds(string(peek))
}

// In case of error just print the error and panic.
func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
