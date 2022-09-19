package pdftypes

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Token definitions
const (
	// Array begin token
	ARRAY_BEGIN = "["
	// Array end token
	ARRAY_END = "]"
	
	// Dictionary begin token
	DICT_BEGIN = "<<"
	// Dictionary end token
	DICT_END = ">>"

	// Hexadecimal begin token
	HEX_BEGIN = "<"
	// Hexadecimal end token
	HEX_END = ">"

	// PdfName begin token
	NAME_BEGIN = "/"

	// PdfReference end token
	REFERENCE_END = "R"

	// String begin token
	STRING_BEGIN = "("
	// String end token
	STRING_END = ")"

	// EOF token
	EOF = "%%EOF"
	// Length of EOF token
	EOFL = len(EOF)

	// Object begin token
	OBJECT = "obj"
	// Object end token
	ENDOBJECT = "endobj"

	// Stream begin token
	STREAM = "stream"
	// Stream end token 
	ENDSTREAM = "endstream"
	// Length of Stream end token
	ENDSTREAML = len(ENDSTREAM)

	// Pdf version prefix token
	PDF_VERSION = "%PDF-"
)

// Check whether the Array begin token appears at the beginning of the line.
func ArrayBegins(str string) bool {
	return strings.HasPrefix(str, ARRAY_BEGIN)
}

// Check whether the Array end token appears at the end of the line.
func ArrayEnds(str string) bool {
	return strings.HasSuffix(str, ARRAY_END)
}

// Check whether the Dictionary begin token appears in the line.
func DictBegins(line_str string) bool {
	return strings.Contains(line_str, DICT_BEGIN)
}

// Check whether the Dictionary end token appears in the line.
func DictEnds(line_str string) bool {
	return strings.Contains(line_str, DICT_END)
}

// Check whether the Hexadecimal begin token appears at the beginning of the line.
func HexBegins(str string) bool {
	return strings.HasPrefix(str, HEX_BEGIN)
}

// Check whether the Hexadecimal end token appears at the end of the line.
func HexEnds(str string) bool {
	return strings.HasSuffix(str, HEX_END)
}

// Check whether the string is a Pdf bool type.
func IsBool(line_str string) bool {
	return line_str == "false" || line_str == "true"
} 

// Check whether the EOF token appears in the line
func IsEOF(line_str string) bool {
	return strings.Contains(line_str, EOF)
}

// Check whether the str is a Pdf null value
func IsNull(line_str string) bool {
	return line_str == "null"
}

// Check whether the Pdf version prefix token appears in the line.
func IsVersion(line_str string) bool {
	return strings.Contains(line_str, PDF_VERSION)
}

// Check whether the number is a Reference.
func IsReference(head string, tokens *[]string) bool {
	if len(*tokens) < 2 {
		return false
	}

	if (*tokens)[1] != REFERENCE_END {
		return false
	}

	_, err1 := strconv.Atoi(head)
	_, err2 := strconv.Atoi((*tokens)[0])
	
	
	return err1 == nil && err2 == nil
}

// Check whether the PdfName begin token appears at the beginning of the line.
func NameBegins(str string) bool {
	return strings.HasPrefix(str, NAME_BEGIN)
}

// Check whether the Object begin token appears in the line.
func ObjectBegins(line_str string) bool {
	return !ObjectEnds(line_str) &&
		strings.HasSuffix(line_str, OBJECT)
}

// Check whether the Object end token appears in the line.
func ObjectEnds(line_str string) bool {
	return strings.HasSuffix(line_str, ENDOBJECT)
}

// Check whether the Stream begin token appears at the beginning of the line.
func StreamBegins(line_str string) bool {
	return strings.HasPrefix(line_str, STREAM)
}

// Check whether the Stream end token appears at the beginning of the line.
func StreamEnds(line_str string) bool {
	return strings.HasPrefix(line_str, ENDSTREAM)
}

// Check whether the String begin token appears at the beginning of the line.
func StringBegins(str string) bool {
	return strings.HasPrefix(str, STRING_BEGIN)
}

// Check whether the String end token appears at the end of the line.
func StringEnds(str string) bool {
	return strings.HasSuffix(str, STRING_END) && !strings.HasSuffix(str, "\\)")
}

// Interface for all basic pdf data types.
// Implement `noOp()` (does nothing) to enable dynamic
// typing and run-time polymorhism.
type PdfDataType interface{
	noOp()
}

// Pdf Name (Symbol/Atom) data type.
type PdfName string
func (n PdfName) noOp() {}

// Pdf number data type.
type PdfNumber float32
func (n PdfNumber) noOp() {}

// Pdf string data type.
type PdfString string
func (s PdfString) noOp() {}

// Pdf reference data type.
type PdfReference struct {
	Object int
	Generation int
}
func (r PdfReference) noOp() {}

// Stringer implementation for PdfReference.
func (r PdfReference) String() string {
	return fmt.Sprintf("%d %d R", r.Object, r.Generation)
}

// Pdf array data type (dynamic).
type PdfArray []PdfDataType
func (a PdfArray) noOp() {}

// Pdf dictionary data type.
type PdfDict map[PdfDataType]PdfDataType
func (d PdfDict) noOp() {}

// Pdf null type
type PdfNull bool
func (n PdfNull) noOp() {}

// Pdf bool type
type PdfBool bool
func (b PdfBool) noOp() {}

// Pdf hexadecimal data type.
type PdfHex string
func (h PdfHex) noOp() {}

// Decode the hexadecimal value.
func (h PdfHex) Decode() ([]byte, error) {
	return hex.DecodeString(string(h))
}
