package pdftypes;

// Token definitions
const DICT_BEGIN = "<<"
const DICT_END = ">>"

const EOF string = "%%EOF"
const EOFL int = len(EOF)

type PdfType interface{ string|int|[]int }

//
type PdfDictEntry[T PdfType] struct {
	key string
	value T
}

//
type PdfDict struct {

}
