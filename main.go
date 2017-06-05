package main

import (
	"io"
	"log"
	"os"
)

// Read everything from input, sort it into natural order, and write
// the result to output.
func NaturalSort(input io.Reader, output io.Writer) error {
	il := &itemList{}
	_, err := io.Copy(il, input)
	if err != nil {
		return err
	}
	il.Flush()
	il.Sort()
	_, err = io.Copy(output, il)
	return err
}

func main() {
	err := NaturalSort(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err.Error())
	}
}
