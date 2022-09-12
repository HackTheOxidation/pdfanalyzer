package pdfobjects

import (
	//"git.magenta.dk/os2datascanner/pdfanalyzer/pdfanalyzer/pdftypes"
)

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

func NewPdfStream(streamtype string, content []byte) *PdfStream {
	return &PdfStream{
		streamtype,
		content,
	}
}

type PdfObject struct {
	pos PdfObjectPosition
	fields []PdfObjectField
	Stream *PdfStream
}

func NewPdfObject() *PdfObject {
	return &PdfObject{
		PdfObjectPosition{x: 0, y: 0},
		make([]PdfObjectField, 1),
		&PdfStream{},
	}
}

