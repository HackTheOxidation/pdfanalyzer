SOURCES=main.cpp PdfSplitter.cpp
CXX=g++
CXXFLAGS = -Wall -std=c++20
EXE=main

all: $(EXE)

$(EXE): $(SOURCES)
	$(CXX) $(CXXFLAGS) -o $@ $^


clean:
	@rm $(EXE)
