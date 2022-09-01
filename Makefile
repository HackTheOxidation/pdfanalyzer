
all: build

build:
	go build -gcflags "-m"

run:
	go build
	./pdfanalyzer

clean:
	go clean
