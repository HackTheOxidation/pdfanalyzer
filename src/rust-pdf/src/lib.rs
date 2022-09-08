mod pdf;

use crate::pdf::{PdfMetadata, PdfObject, PdfSource};


pub struct PdfAST {
    source: Box<PdfSource>,
    metadata: PdfMetadata,
    objects: Vec<PdfObject>,
}

impl PdfAST {
    pub fn new(filename: String) -> PdfAST {
	let source = PdfSource::new(filename);
	
	PdfAST {
	    source: Box::new(source),
	    metadata: PdfMetadata::extract_metadata(source),
	    objects: PdfAST::parse(source),
	}
    }

    fn parse(source: PdfSource) -> Vec<PdfObject> {
	vec![PdfObject{}]
    }
}

