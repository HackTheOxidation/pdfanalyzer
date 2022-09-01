package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfobjects"
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdftypes"
)

// PdfParser interface
//
// Defines a common interface for what it means to be a PdfParser.
// This allows for different implementations of parsers with different
// strategies.
type PdfParser interface {
	parseName() pdftypes.PdfName
	parseString() pdftypes.PdfString
	parseReference() pdftypes.PdfReference
	parseNumber() pdftypes.PdfNumber
	parseArray() pdftypes.PdfArray
	parseDictionary(line_str string) pdftypes.PdfDict
}

// Read an object from the file. Panics if EOF is reached.
func (r *PdfReader) readObject(line_number int) *pdfobjects.PdfObject {
	var err error = nil;
	var line []byte;
	pobj := pdfobjects.NewPdfObject()

	for !pdftypes.ObjectEnds(string(line)) {
		line, _, err = r.reader.ReadLine(); 
		check(err)
		line_number++

		r.dispatch(line, pobj, line_number)
	}

	return pobj
}

// Read a stream from the file and insert it into `pobj`.
// Panics if EOF is reached.
func (r *PdfReader) readStream(pobj *pdfobjects.PdfObject, line_number int) {
	var buffer []byte;
	
	for !r.reader.IsEndstream() {
		b, err := r.reader.ReadByte()
		parserError(err, line_number)
		buffer = append(buffer, b)
	}

	pobj.Stream = pdfobjects.NewPdfStream("Stream", buffer)
}

// Decide what construct to read next. Either a stream or a dictionary.
func (r *PdfReader) dispatch(line []byte, pobj *pdfobjects.PdfObject, line_number int) {
	line_str := string(line)
	if pdftypes.StreamBegins(line_str) {
		r.readStream(pobj, line_number)
	} else if pdftypes.DictBegins(line_str) {
		tokens := r.tokenizeDict(line_str, line_number)
		result, _ := r.parseValue(&tokens, line_number)

		dict := result.(pdftypes.PdfDict)
		
		pobj.SetDict(dict)
	} 
}

func formatLine(line_str string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						line_str,
						">>", " >> "),
					"]", " ] "),
				"[", " [ "),
			"/", " /"),
		"<<", " << ")
}

func (r *PdfReader) tokenizeDict(line_str string, line_number int) []string {
	buffer := ""

	buffer += formatLine(line_str)

	bal_dict := strings.Count(buffer, pdftypes.DICT_BEGIN) - strings.Count(buffer, pdftypes.DICT_END)

	for bal_dict != 0 {
		line, _, err := r.reader.ReadLine()
		line_number++
		parserError(err, line_number)
		line_str = formatLine(string(line))

		buffer += " " + line_str
		
		bal_dict += strings.Count(line_str, pdftypes.DICT_BEGIN)
		bal_dict -= strings.Count(line_str, pdftypes.DICT_END)
	}

	tokens := strings.Split(buffer, " ")

	return tokens
}

// Read and parse a pdf dictionary from the file and insert is into `pobj`.
func (r *PdfReader) parseDictionary(tokens *[]string, line_number int) pdftypes.PdfDict {
	dict := make(pdftypes.PdfDict)

	cursor := strings.TrimSpace((*tokens)[0])

	for !pdftypes.DictEnds(cursor) {		
		key, err := r.parseValue(tokens, line_number)
		
		if err != nil {
			break
		}
		
		value, err := r.parseValue(tokens, line_number)

		dict[key] = value
	}

	return dict
}

// Dispatch and parse a value of an appropriate PdfDataType
func (r *PdfReader) parseValue(tokens *[]string, line_number int) (pdftypes.PdfDataType, error) {
	if len(*tokens) == 0 {
		return nil, errors.New("Out of tokens")
	}
	
	head := strings.TrimSpace((*tokens)[0])

	*tokens = (*tokens)[1:]

	if head == "" {
		return r.parseValue(tokens, line_number)
	}

	if head == pdftypes.DICT_END {
		return nil, errors.New("End of dictionary")
	}

	if head == pdftypes.ARRAY_END {
		return nil, errors.New("End of array")
	}
	
	if pdftypes.StringBegins(head) {
		return r.parseString(head, tokens, line_number), nil
	} else if pdftypes.DictBegins(head) {
		return r.parseDictionary(tokens, line_number), nil		
	} else if pdftypes.ArrayBegins(head) {
		return r.parseArray(tokens, line_number), nil	
	} else if pdftypes.HexBegins(head) {
		return r.parseHex(strings.TrimSpace(head), line_number), nil
	} else if pdftypes.NameBegins(head) {
		return r.parseName(strings.TrimSpace(head), line_number), nil
	} else if pdftypes.IsNull(head) {
		return pdftypes.PdfNull(false), nil
	} else if pdftypes.IsBool(head) {
		return r.parseBool(head), nil
	} else {
		if pdftypes.IsReference(head, tokens) {
			return r.parseReference(head, tokens, line_number), nil
		} else {
			return r.parseNumber(head, line_number), nil
		}
	}
}

func (r *PdfReader) parseBool(line_str string) pdftypes.PdfBool {
	if line_str == "true" {
		return pdftypes.PdfBool(true)
	} else {
		return pdftypes.PdfBool(false)
	}
}

// Parse a value of type: PdfName
func (r *PdfReader) parseName(str string, line_number int) pdftypes.PdfName {
	trimmed := strings.TrimPrefix(str, " ")
	if !strings.HasPrefix(trimmed, "/") {
		unexpectedToken("/", str, line_number)
	}
	return pdftypes.PdfName(trimmed)
}

// Parse a value of type: PdfString
func (r *PdfReader) parseString(head string, tokens *[]string, line_number int) pdftypes.PdfString {
	buffer := head

	balance := strings.Count(buffer, pdftypes.STRING_BEGIN) - strings.Count(buffer, pdftypes.STRING_END)
	if balance != 0 {
		cursor := strings.TrimSpace((*tokens)[0])
		buffer += cursor + " "
	
		*tokens = (*tokens)[1:]

		return r.parseString(buffer, tokens, line_number)
	}

	return pdftypes.PdfString(buffer)	
}

// Parse a value of type: PdfArray
func (r *PdfReader) parseArray(tokens *[]string, line_number int) pdftypes.PdfArray {
	array := make(pdftypes.PdfArray, 0)
	
	cursor := strings.TrimSpace((*tokens)[0])

	for !pdftypes.ArrayEnds(cursor) {
		element, err := r.parseValue(tokens, line_number)

		if err != nil {
			break
		}
		
		array = append(array, element)

		cursor = strings.TrimSpace((*tokens)[0])
	}

	*tokens = (*tokens)[1:]

	return array
}

// Parse a value of type: PdfHex
func (r *PdfReader) parseHex(str string, line_number int) pdftypes.PdfHex {
	if !pdftypes.HexEnds(str) {
		missingDelimiter(pdftypes.HEX_END, line_number)
	}

	return pdftypes.PdfHex(
		strings.TrimSuffix(
			strings.TrimPrefix(str, pdftypes.HEX_BEGIN),
			pdftypes.HEX_END,
		),
	)
}

// Parse a value of type: PdfNumber
func (r *PdfReader) parseNumber(str string, line_number int) pdftypes.PdfNumber {
	num, err := strconv.ParseFloat(str, 32)
	parserError(err, line_number)
	return pdftypes.PdfNumber(num)
}

// Parse a value of type: PdfReference
func (r *PdfReader) parseReference(head string, tokens *[]string, line_number int) pdftypes.PdfReference {
	
	object, err := strconv.Atoi(head)
	parserError(err, line_number)
	
	generation, err := strconv.Atoi((*tokens)[0])
	parserError(err, line_number)

	*tokens = (*tokens)[2:]

	return pdftypes.PdfReference{
		Object: object,
		Generation: generation,
	}
}

// Terminating error handler for ParserError - Unexpected Token.
func unexpectedToken(expected string, actual string, line_number int) {
	panic(fmt.Sprintf("ParserError - Unexpected Token at line %d: got %s, expected %s\n", line_number, actual, expected))
}

// Terminating error handler for ParserError - Missing delimiter.
func missingDelimiter(missing string, line_number int) {
	panic(fmt.Sprintf("ParserError - Missing delimiter at line %d: %v\n", line_number, missing))
}

func parserError(err error, line_number int) {
	if err != nil {
		panic(fmt.Sprintf("Error at line %d: %v", line_number, err))
	}
}
