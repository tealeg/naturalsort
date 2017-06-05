package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var inputType string
var inputFile string
var outputType string
var outputFile string

func init() {
	setupFlags(&inputType, &inputFile, &outputType, &outputFile)
}

func setupFlags(inputType, inputFile, outputType, outputFile *string) {
	const (
		defaultInputType  = "stdin"
		inputTypeUsage    = "The type of input to read: (stdin | file | gzip)"
		defaultInputFile  = ""
		inputFileUsage    = "A path to a file to read input from when input_type is either file or gzip."
		defaultOutputType = "stdout"
		outputTypeUsage   = "The type of output to write: (stdout | file | gzip)"
		defaultOutputFile = ""
		outputFileUsage   = "A path to a file to write output to when output_type is either file or gzip."
	)
	flag.StringVar(inputType, "it", defaultInputType, inputTypeUsage)
	flag.StringVar(inputFile, "if", defaultInputFile, inputFileUsage)
	flag.StringVar(outputType, "ot", defaultOutputType, outputTypeUsage)
	flag.StringVar(outputFile, "of", defaultOutputFile, outputFileUsage)

}

// Read everything from input, sort it into natural order, and write
// the result to output.
func NaturalSort(input io.Reader, output io.Writer) error {
	il := &itemList{}
	_, err := io.Copy(il, input)
	if err != nil {
		return err
	}
	err = il.Flush()
	if err != nil {
		return err
	}
	il.Sort()
	_, err = io.Copy(output, il)
	return err
}

// Return the io.ReadCloser implentation implied by inputType (and, where necessary, path).
func getInput(inputType, path string) (reader io.ReadCloser, err error) {
	switch strings.ToLower(inputType) {
	case "stdin":
		reader = os.Stdin
	case "file":
		if path == "" {
			err = fmt.Errorf("When input type is 'file', a path must be provided.")
			return
		}
		reader, err = os.Open(path)
	case "gzip":
		if path == "" {
			err = fmt.Errorf("When input type is 'gzip', a path must be provided.")
			return
		}
		var freader io.Reader
		freader, err = os.Open(path)
		if err == nil {
			reader, err = gzip.NewReader(freader)
		}
	}
	return
}

// getOutput takes and outputType and a path (optionally empty string)
// and returns an io.WriteCloser implementation.
func getOutput(outputType, path string) (writer io.WriteCloser, err error) {
	switch strings.ToLower(outputType) {
	case "stdout":
		writer = os.Stdout
	case "file":
		if path == "" {
			err = fmt.Errorf("When output type is 'file', a path must be provided.")
			return
		}
		writer, err = os.Create(path)
	case "gzip":
		if path == "" {
			err = fmt.Errorf("When output type is 'gzip', a path must be provided.")
			return
		}
		var fwrite io.Writer
		fwrite, err = os.Create(path)
		if err == nil {
			writer = gzip.NewWriter(fwrite)
		}
	}
	return
}

func main() {
	flag.Parse()
	inputReader, err := getInput(inputType, inputFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer inputReader.Close()

	outputWriter, err := getOutput(outputType, outputFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer outputWriter.Close()

	err = NaturalSort(inputReader, outputWriter)
	if err != nil {
		log.Fatal(err.Error())
	}
}
