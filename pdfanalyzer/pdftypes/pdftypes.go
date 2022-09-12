package pdftypes

// Token definitions
const (
	DICT_BEGIN = "<<"
	DICT_END = ">>"

	EOF string = "%%EOF"
	EOFL int = len(EOF)

	ENDSTREAM string = "endstream"
	ENDSTREAML int = len(ENDSTREAM)
)


type PdfType interface{ string|int|[]int }

//
type PdfDictEntry[T PdfType] struct {
	key string
	value T
}

//
type PdfDict struct {

}
