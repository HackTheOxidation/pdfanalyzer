package pdfobjects

import (
	"bytes"
	"compress/lzw"
	"compress/zlib"
	"encoding/ascii85"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"git.magenta.dk/os2datascanner/pdfanalyzer/pdftypes"
)

// Wrapper for pdf stream.
type PdfStream struct {
	streamtype string
	content []byte
}

// Create a new `PdfStream` from byte array and streamtype.
func NewPdfStream(streamtype string, content []byte) PdfStream {
	return PdfStream{
		streamtype,
		content,
	}
}

// Extract the contents of a stream.
func (s PdfStream) Extract(pobj *PdfObject, cmap *CMap) ([]byte, bool, error) {
	// If the stream is empty, don't return anything.
	if len(s.content) == 0 {
		return nil, false, errors.New("This stream is empty.")
	}

	// If the stream contains an image, don't return anything.
	// TODO: Inline OCR text extraction.
	if pobj.IsImage() {
		return nil, false, nil
	}

	// Pass the stream contents to an appropriate extraction handler.
	switch pobj.GetEncoding() {
	case pdftypes.ASCII85DECODE:
		return extractASCII85(s.content, cmap)

	case pdftypes.ASCIIHEXDECODE:
		return extractASCIIHEX(s.content, cmap)

	case pdftypes.LZWDECODE:
		return extractLZW(s.content, cmap)

	case pdftypes.FLATEDECODE:
		return extractZlib(s.content, cmap)
		
	default:
		return extractStrings(s.content, cmap)
	}
}

func extractASCII85(content []byte, cmap *CMap) ([]byte, bool, error) {
	b := bytes.NewReader(content)

	extracted, err := io.ReadAll(ascii85.NewDecoder(b))
	if err != nil {
		return nil, false, err
	}
	
	return extractStrings(extracted, cmap)
}

func extractASCIIHEX(content []byte, cmap *CMap) ([]byte, bool, error) {
	extracted, _ := hex.DecodeString(string(content))

	return extractStrings(extracted, cmap)
}

func extractLZW(content []byte, cmap *CMap) ([]byte, bool, error) {
	b := bytes.NewReader(content)

	extracted, err := io.ReadAll(lzw.NewReader(b, lzw.LSB, 256))
	if err != nil {
		return nil, false, err
	}
	
	return extractStrings(extracted, cmap)
}

// Decodes with zlib and extracts the textual contents from a stream.
func extractZlib(content []byte, cmap *CMap) ([]byte, bool, error) {
	b := bytes.NewReader(content)
	
	rc, err := zlib.NewReader(b)
	if err != nil {
		return nil, false, err
	}
	
	extracted, err := io.ReadAll(rc)
	if err != nil {
		return nil, false, err
	}
	
	return extractStrings(extracted, cmap)
}

// Extract in-stream text from strings
func extractStrings(content []byte, cmap *CMap) ([]byte, bool, error) {
	// Search for a CMap definition and try to parse it.
	parseCMap(bytes.NewReader(content), cmap)

	// Search for Text as strings or byte mappings.
	reader := bytes.NewReader(content)
	text := make([]byte, 0)
	process := false

	// Flags for changing state.
	in_text, in_str, in_hex := false, false, false
	// Previous byte. 
	var prev byte = 0

	// Read the stream content byte by byte
	for b, err := reader.ReadByte(); err == nil; b, err = reader.ReadByte() {
		if in_text {
			if in_str { // String parsing state.
				if b == ')' && prev != '\\' {
					// If an unescaped ')' is encountered the string is terminated.
					in_str = false
				} else {
					// Otherwise just parse the byte.
					text = append(text, b)
				}
			} else if in_hex { // Hexadecimal parsing state.
				if b == '>' {
					// If an '>' is encountered the hex value is terminated.
					in_hex = false
					text = append(text, b)
				} else {
					// Otherwise just parse the byte.
					text = append(text, b)
				}
			} else { // Search for either a string or a hex value.
				if b == '(' {
					// A string has begun. Start parsing it.
					in_str = true
				} else if b == '*' && prev == 'T' {
					// T* is the newline token.
					text = append(text, '\n')
				} else if b == 'T' && prev == 'E' {
					// ET signifies end of text. Search for a new text block.
					in_text = false
				} else if b == '-' {
					// Add space to text as '-' signifies a 'SPC'-character. 
					text = append(text, ' ')
				} else if b == '<' {
					// A hexadecimal value is encountered. Start parsing it.
					in_hex = true
					process = true
					text = append(text, b)
				}
			}
		} else {
			// If BT is encountered. Start searching for strings.
			if b == 'T' && prev == 'B' {
				in_text = true
			}
		}

		// Remember the previous byte.
		prev = b
	}
	
	return text, process, nil
}

func checkDecoding(err error, hex_str string) {
	if err != nil {
		fmt.Printf("Error - Unable to decode hexadecimal value: <%v>\n", hex_str)
		os.Exit(-1)
	}
}
