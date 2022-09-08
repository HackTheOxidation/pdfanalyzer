package main;

import (
	"bufio"
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	fp, err := os.Open("file.txt")
	check(err)

	reader := bufio.NewReader(fp)
	fmt.Printf("Bytes remaining: %d\n", reader.Buffered())

	b, err := reader.ReadByte()
	check(err)
	fmt.Printf("Read byte: %c, Bytes remaining: %d\n", b, reader.Buffered())

	b, err = reader.ReadByte()
	check(err)
	fmt.Printf("Read byte: %c, Bytes remaining: %d\n", b, reader.Buffered())

	reader = bufio.NewReader(fp)
	b, err = reader.ReadByte()
	check(err)
	fmt.Printf("Read byte: %c, Bytes remaining: %d\n", b, reader.Buffered())

	b, err = reader.ReadByte()
	check(err)
	fmt.Printf("Read byte: %c, Bytes remaining: %d\n", b, reader.Buffered())

	fp.Close()
}
