package main;

import (
	"bytes"
	"bufio"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const EOF string = "%%EOF"
const EOFL int = len(EOF)

const ENDSTREAM string = "endstream"
const ENDSTREAML int = len(ENDSTREAM)

type RefreshingReader struct {
	fp *os.File
	reader *bufio.Reader
}

func NewRefreshingReader(filename string) (*RefreshingReader, error) {
	fp, err := os.Open(filename)
	reader := bufio.NewReader(fp)
	return &RefreshingReader{
		fp,
		reader,
	}, err
}

func (r *RefreshingReader) reloadBuffer() {
	r.reader = bufio.NewReader(r.fp)
}

func (r *RefreshingReader) Close() {
	r.fp.Close()
}

func (r *RefreshingReader) Name() string {
	return r.fp.Name()
}

func (r *RefreshingReader) ReadByte() (byte, error) {
	b, err := r.reader.ReadByte()
	if err == io.EOF {
		r.reloadBuffer()
		return r.ReadByte()
	}
	return b, err
}

func (r *RefreshingReader) Peek(n int) ([]byte, error) {
	return r.reader.Peek(n)
}

func (r *RefreshingReader) ReadLine() ([]byte, bool, error) {
	return r.reader.ReadLine()
}

func (r *RefreshingReader) IsEOF() bool {
	peek, err := r.reader.Peek(EOFL)
	check(err)
	return string(peek) == EOF
}

func (r *RefreshingReader) IsEndstream() bool {
	peek, err := r.Peek(ENDSTREAML)
	check(err)
	return string(peek) == ENDSTREAM
}

type PdfAST struct {
	objects []PdfObject
	metadata PdfMetadata
}

type PdfMetadata struct {
	version string
}

type PdfObjectField struct {
	field string
	value string 
}

type PdfObjectPosition struct {
	x int
	y int
}

type PdfStream struct {
	streamtype string
	content []byte
}

type PdfObject struct {
	pos PdfObjectPosition
	fields []PdfObjectField
	stream PdfStream
}

func NewPdfObject() *PdfObject {
	return &PdfObject{
		PdfObjectPosition{x: 0, y: 0},
		make([]PdfObjectField, 1),
		PdfStream{},
	}
}


func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func ReadObject(reader *RefreshingReader) *PdfObject {
	var err error = nil;
	var line []byte;
	pobj := NewPdfObject()

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

func readStream(reader *RefreshingReader, pobj *PdfObject) {
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
	pobj.stream = PdfStream{ "Stream", stream }
	//fmt.Printf("stream: %s\n", stream)
}

func readDict(reader *RefreshingReader, pobj *PdfObject) {
	
}

func dispatch(line []byte, reader *RefreshingReader, pobj *PdfObject) {
	line_str := string(line)
	if strings.HasPrefix(line_str, "stream") {
		readStream(reader, pobj)
	} else if strings.HasPrefix(line_str, "<<") {
		readDict(reader, pobj)
	} 
}

func main() {
	file_name := "../../assets/gcc.pdf"
	reader, err := NewRefreshingReader(file_name)
	check(err)
	count := 0

	for !reader.IsEOF() {
		ReadObject(reader)
		count++
		fmt.Printf("Read Object #%d.\n", count)
	}
	
	// fp, err := os.Open("../../assets/gcc.pdf")
	// check(err)

	// name := fp.Name()
	// fmt.Printf("File name: %s\n", name)

	// reader := bufio.NewReader(fp)

	// peek, err := reader.Peek(1)
	// r := reader.Buffered()
	// check(err)
	// fmt.Printf("Buffered: %d, Peak: %x\n", r, peek)
	// for r := reader.Buffered(); r > 0; r = reader.Buffered() {
	// 	ReadObject(reader)
	// }
	
	// fp.Close()
}
