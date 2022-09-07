#include "PdfSplitter.hpp"


PdfSplitter::PdfSplitter(std::string file_name) {
  pdf_stream = std::make_shared<std::ifstream>(file_name);
}


PdfSplitter::~PdfSplitter() {
  refs--;

  if (refs < 1)
    pdf_stream->close();
}

PdfSplitter::PdfSplitter(const PdfSplitter & other)
  : pdf_stream(other.pdf_stream), refs(other.refs + 1) {}

PdfSplitter& PdfSplitter::operator=(const PdfSplitter & other) {
  this->pdf_stream = other.pdf_stream;
  this->refs = other.refs;

  return *this;
}
