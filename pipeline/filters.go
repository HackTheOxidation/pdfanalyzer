package pipeline

import (
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdfobjects"
	"git.magenta.dk/os2datascanner/pdfanalyzer/pdftypes"
)

// Type signature for functions for the filtering step.
type FilterFunction func (in <-chan pdfobjects.PdfObject, out chan<- pdfobjects.PdfObject)

// Runs the filtering stage with the filter function `f`.
func runFilterStage(f FilterFunction, in <-chan pdfobjects.PdfObject, out chan<- pdfobjects.PdfObject) {
	f(in, out)
}

// Base function for filtering.
func filter(in <-chan pdfobjects.PdfObject, out chan<- pdfobjects.PdfObject, condition func (pdfobjects.PdfObject) bool) {
	defer close(out)
	for data := range in {
		if condition(data) {
			out <- data
		}
	}
}

// The identity filter is just like the identity function: it does nothing.
func IdentityFilter(in <-chan pdfobjects.PdfObject, out chan<- pdfobjects.PdfObject) {
	filter(in, out, func (_ pdfobjects.PdfObject) bool { return true })
}

// Removes all objects that doesn't contain any text.
func TextOnlyFilter(in <-chan pdfobjects.PdfObject, out chan<- pdfobjects.PdfObject) {
	filter(in, out, func (po pdfobjects.PdfObject) bool { return po.IsText() })
}

// Removes all objects expect images.
func ImageOnlyFilter(in <-chan pdfobjects.PdfObject, out chan<- pdfobjects.PdfObject) {
	filter(in, out, func(po pdfobjects.PdfObject) bool { return po.IsImage() })
}

// Removes all objects expect objects of type /ObjStm
func ObjStmOnlyFilter(in <-chan pdfobjects.PdfObject, out chan<- pdfobjects.PdfObject) {
	filter(in, out, func(po pdfobjects.PdfObject) bool { return po.GetType() == pdftypes.OBJSTM })
}
