package main

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

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
		reader, err = os.Open(path)
	case "gzip":
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
		writer, err = os.Create(path)
	case "gzip":
		var fwrite io.Writer
		fwrite, err = os.Create(path)
		if err == nil {
			writer = gzip.NewWriter(fwrite)
		}
	}
	return
}

func main() {
	var inputType string
	var inputFile string
	var outputType string
	var outputFile string

	setupFlags(&inputType, &inputFile, &outputType, &outputFile)
	flag.Parse()
	inputReader, err := getInput(inputType, inputFile)
	defer inputReader.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	outputWriter, err := getOutput(outputType, outputFile)
	defer outputWriter.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = NaturalSort(inputReader, outputWriter)
	if err != nil {
		log.Fatal(err.Error())
	}
}
