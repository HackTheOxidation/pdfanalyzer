#include <iostream>
#include <memory>
#include <fstream>
#include <cstddef>

#ifndef PDFSPLITTER_HPP_
#define PDFSPLITTER_HPP_

class PdfSplitter {
public:
  PdfSplitter(std::string);
  PdfSplitter(const PdfSplitter &);
  ~PdfSplitter();
  PdfSplitter& operator=(const PdfSplitter &);

private:
  std::shared_ptr<std::ifstream> pdf_stream;
  std::size_t refs = 1;
};

#endif
