package pdftypes

const (
	// Object type identifiers
	OBJ_TYPE PdfName = "/Type"
	XOBJECT PdfName = "/XObject"
	OBJSTM PdfName = "/ObjStm"
	PAGE PdfName = "/Page"
	PAGES PdfName = "/Pages"

	// Compression specifier
	FILTER PdfName = "/Filter"

	// Compression Methods
	ASCIIHEXDECODE PdfName = "/ASCIIHexDecode"
	ASCII85DECODE PdfName = "/ASCII85Decode"
	LZWDECODE PdfName = "/LZWDecode"
	FLATEDECODE PdfName = "/FlateDecode"
	RUNLENGTHDECODE PdfName = "/RunLengthDecode"
	CCITTFAXDECODE PdfName = "/CCITTFaxDecode"
	JBIG2DECODE PdfName = "/JBIG2Decode"
	DCTDECODE PdfName = "/DCTDecode"
	JPXDECODE PdfName = "/JPXDecode"
)
