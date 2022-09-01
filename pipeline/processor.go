package pipeline

import (
	"bytes"
	"encoding/hex"

	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfobjects"
)

// Type signature for functions for the processor step.
type ProcessorFunction func (in <-chan ExtractorResult, out chan<- ProcessorResult, cmap *pdfobjects.CMap)

func runProcessorStage(f ProcessorFunction, in <-chan ExtractorResult, out chan<- ProcessorResult, cmap *pdfobjects.CMap) {
	f(in, out, cmap)
}

type ProcessorResult struct {
	stream string
	err error
}

func NewProcessorResult(stream string, err error) ProcessorResult {
	return ProcessorResult{
		stream,
		err,
	}
}

func transform(hex_buffer []byte, cmap *pdfobjects.CMap) string {
	if len(hex_buffer) == 0 {
		return ""
	}
	
	if len(hex_buffer) < 4 {
		res, err := cmap.Convert(string(hex_buffer))
		if err != nil {
			return ""
		} else {
			return res
		}
	} else {
		head, tail := hex_buffer[:4], hex_buffer[4:]
		res, err := cmap.Convert(string(head))
		if err != nil {
			return "" + transform(tail, cmap)
		} else {
			return res + transform(tail, cmap)
		}
	}
}

func mapCharacters(extracted ExtractorResult, cmap *pdfobjects.CMap) ProcessorResult {
	text, hex_buffer := make([]byte, 0), make([]byte, 0)
	reader := bytes.NewReader([]byte(extracted.stream))

	in_hex := false

	for b, err := reader.ReadByte(); err == nil; b, err = reader.ReadByte() {
		if in_hex {
			if b == '>' {
				// '>'
				in_hex = false

				// Transform and decode hex.
				decoded, _ := hex.DecodeString(transform(hex_buffer, cmap))
				text = append(text, decoded...)
				
				// Reset the hex buffer
				hex_buffer = make([]byte, 0)
			} else {
				// Otherwise just append the byte to the hex buffer.
				hex_buffer = append(hex_buffer, b)
			}
		} else {
			if b == '<' {
				// '<'
				in_hex = true
			} else if b == ' ' {
				text = append(text, b)
			}
		}
	}

	return NewProcessorResult(string(text), nil)
}

func CMapProcessor(in <-chan ExtractorResult, out chan<- ProcessorResult, cmap *pdfobjects.CMap) {
	defer close(out)
	for data := range in {
		if data.process {
			out <- mapCharacters(data, cmap)
		} else {
			out <- data.ToProcessorResult()
		}
	}
}
