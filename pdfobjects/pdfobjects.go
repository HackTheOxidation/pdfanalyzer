package pdfobjects

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"

	"git.magenta.dk/os2datascanner/pdfanalyzer/pdftypes"
)

// Wrapper for read pdf file.
type Pdf struct {
	name string
	version string
	objects []*PdfObject
	count int
}

// Create a new empty `Pdf` struct.
func NewPdf(name string) *Pdf {
	return &Pdf{
		name,
		"",
		make([]*PdfObject, 0),
		0,
	}
}

// Get the name of the pdf file.
func (pdf Pdf) Name() string {
	return pdf.name
}

// Update the version of the file.
func (pdf *Pdf) SetVersion(version string) {
	pdf.version = version
}

// Retrieve the version of the file.
func (pdf Pdf) Version() string {
	return pdf.version
}

// Append a `PdfObject` to the wrappers internal list of objects.
func (pdf *Pdf) AppendObject(obj *PdfObject) {
	pdf.objects = append(pdf.objects, obj)
	pdf.count++
}

// Return the number of read objects.
func (pdf Pdf) Count() int {
	return pdf.count
}

// Return object with index `i`.
func (pdf Pdf) GetObject(i int) (*PdfObject, error) {
	if i >= pdf.Count() {
		return nil, errors.New("Index out of Bounds.")
	} else {
		return pdf.objects[i], nil
	}
}

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
func (s PdfStream) Extract(pobj *PdfObject) ([]byte, error) {
	if len(s.content) == 0 {
		return nil, errors.New("This stream is empty.")
	}
	
	b := bytes.NewReader(s.content)
	rc, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	
	return io.ReadAll(rc)
}

// Wrapper for pdf object.
type PdfObject struct {
	pos pdftypes.PdfReference
	dict pdftypes.PdfDict
	Stream PdfStream
}

// Create a new empty `PdfObject`.
func NewPdfObject() *PdfObject {
	return &PdfObject{
		pdftypes.PdfReference{Object: 0, Generation: 0},
		make(pdftypes.PdfDict, 0),
		PdfStream{},
	}
}

// Update the dictionary associated with the PdfObject.
func (pobj *PdfObject) SetDict(dict pdftypes.PdfDict) {
	pobj.dict = dict
}

// Helper function for extracting the stream of the PdfObject.
func (pobj *PdfObject) ExtractStream() ([]byte, error) {
	return pobj.Stream.Extract(pobj)
}
