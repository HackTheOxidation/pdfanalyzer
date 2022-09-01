package pdfobjects

import (
	"errors"

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
func (pobj *PdfObject) ExtractStream(cmap *CMap) (string, bool, error) {
	stream, process, err := pobj.Stream.Extract(pobj, cmap)
	
	if stream == nil {
		return "", process, err
	}
	
	return string(stream), process, err
}

// Returns the type of the object.
func (pobj PdfObject) GetType() pdftypes.PdfName {
	t, ok := pobj.dict[pdftypes.OBJ_TYPE].(pdftypes.PdfName)

	if ok {
		return t
	} else {
		return pdftypes.PdfName("")
	}
}

// Checks if the object is an image
func (pobj PdfObject) IsImage() bool {
	if pobj.dict[pdftypes.OBJ_TYPE] == pdftypes.XOBJECT {
		return true
	}

	return false
}

// Checks if the object is just text.
func (pobj PdfObject) IsText() bool {
	return pobj.dict[pdftypes.OBJ_TYPE] == nil
}

// Check whether the stream of the object is encoded.
func (pobj PdfObject) IsEncoded() bool {
	return pobj.dict[pdftypes.FILTER] != nil
}

// Get the encoding type of the associated stream.
func (pobj PdfObject) GetEncoding() pdftypes.PdfName {
	if !pobj.IsEncoded() {
		return ""
	}

	return pobj.dict[pdftypes.FILTER].(pdftypes.PdfName) 
}
