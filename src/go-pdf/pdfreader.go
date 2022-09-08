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
	PdfObjectPosition
	fields []PdfObjectField
	PdfStream
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
		panic(err)
	}
}

func ReadObject(reader *bufio.Reader) {
	var err error = nil;
	var line []byte;
	pobj := NewPdfObject()

	for !strings.Contains(string(line), "endobj") {
		line, _, err = reader.ReadLine(); 
		check(err)
		fmt.Printf("%s\n", line)

		dispatch(line, reader, pobj)
	}
}

func splitKeyValue(words []string) (string, string, error) {
	if len(words) == 2 {
		return words[0], words[1], nil
	} else {
		return "", "", errors.New("")
	}
}

func readStream(reader *bufio.Reader, pobj *PdfObject) {
	var buffer []byte;
	peek, err := reader.Peek(len("endstream"))
	check(err)
	for string(peek) != "endstream" {
		b, err := reader.ReadByte()
		check(err)
		buffer = append(buffer, b)
	}

	b := bytes.NewReader(buffer)
	rc, err := zlib.NewReader(b)
	check(err)

	stream, err := io.ReadAll(rc) 
	check(err)
	fmt.Printf("stream: %s\n", stream)
}

func readDict(reader *bufio.Reader, pobj *PdfObject) {

}

func dispatch(line []byte, reader *bufio.Reader, pobj *PdfObject) {
	line_str := string(line)
	if strings.HasPrefix(line_str, "stream") {
		readStream(reader, pobj)
	} else if strings.HasPrefix(line_str, "<<") {
		readDict(reader, pobj)
	} 
}

func main() {
	fp, err := os.Open("../../assets/gcc.pdf")
	check(err)

	name := fp.Name()
	fmt.Printf("File name: %s\n", name)

	reader := bufio.NewReader(fp)
	peek, err := reader.Peek(EOFL)
	check(err)
	for string(peek) != EOF {
		ReadObject(reader)

		peek, err = reader.Peek(EOFL)
		check(err)
	}
	
	fp.Close()
}
