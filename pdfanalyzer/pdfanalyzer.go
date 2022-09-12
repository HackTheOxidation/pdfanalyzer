package pdfanalyzer;

import (
	"bytes"
	"bufio"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfanalyzer/pdftypes"
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfanalyzer/pdfobjects"
)

type PdfReader struct {
	filename string
	reader *refreshingReader
}

func NewPdfReader(filename string) (*PdfReader, error) {
	reader, err := newRefreshingReader(filename)
	return &PdfReader{
		filename,
		reader,
	}, err
}

func (r *PdfReader) ReadAll() {
	count := 0

	for !r.reader.IsEOF() {
		r.reader.readObject()
		count++
		fmt.Printf("Read Object #%d.\n", count)
	}
	
	r.reader.Close()
}

type refreshingReader struct {
	fp *os.File
	reader *bufio.Reader
}

func newRefreshingReader(filename string) (*refreshingReader, error) {
	fp, err := os.Open(filename)
	reader := bufio.NewReader(fp)
	return &refreshingReader{
		fp,
		reader,
	}, err
}

func (r *refreshingReader) reloadBuffer() {
	r.reader = bufio.NewReader(r.fp)
}

func (r *refreshingReader) Close() {
	r.fp.Close()
}

func (r *refreshingReader) Name() string {
	return r.fp.Name()
}

func (r *refreshingReader) ReadByte() (byte, error) {
	b, err := r.reader.ReadByte()
	if err == io.EOF {
		r.reloadBuffer()
		return r.ReadByte()
	}
	return b, err
}

func (r *refreshingReader) Peek(n int) ([]byte, error) {
	return r.reader.Peek(n)
}

func (r *refreshingReader) ReadLine() ([]byte, bool, error) {
	return r.reader.ReadLine()
}

func (r *refreshingReader) IsEOF() bool {
	peek, err := r.reader.Peek(pdftypes.EOFL)
	check(err)
	return string(peek) == pdftypes.EOF
}

func (r *refreshingReader) IsEndstream() bool {
	peek, err := r.Peek(pdftypes.ENDSTREAML)
	check(err)
	return string(peek) == pdftypes.ENDSTREAM
}




func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (reader *refreshingReader) readObject() *pdfobjects.PdfObject {
	var err error = nil;
	var line []byte;
	pobj := pdfobjects.NewPdfObject()

	for !strings.Contains(string(line), "endobj") {
		line, _, err = reader.ReadLine(); 
		check(err)
		//fmt.Printf("%s\n", line)

		dispatch(line, reader, pobj)
	}

	return pobj
}

func splitKeyValue(words []string) (string, string, error) {
	if len(words) == 2 {
		return words[0], words[1], nil
	} else {
		return "", "", errors.New("")
	}
}

func readStream(reader *refreshingReader, pobj *pdfobjects.PdfObject) {
	var buffer []byte;
	
	for !reader.IsEndstream() {
		b, err := reader.ReadByte()
		check(err)
		buffer = append(buffer, b)
	}

	b := bytes.NewReader(buffer)
	rc, err := zlib.NewReader(b)
	check(err)

	stream, err := io.ReadAll(rc) 
	check(err)
	pobj.Stream = pdfobjects.NewPdfStream("Stream", stream)
	//fmt.Printf("stream: %s\n", stream)
}

func readDict(reader *refreshingReader, pobj *pdfobjects.PdfObject) {
	
}

func dispatch(line []byte, reader *refreshingReader, pobj *pdfobjects.PdfObject) {
	line_str := string(line)
	if strings.HasPrefix(line_str, "stream") {
		readStream(reader, pobj)
	} else if strings.HasPrefix(line_str, "<<") {
		readDict(reader, pobj)
	} 
}

