package main

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
}

// getInput takes an input type and a path (optionally empty string)
// and returns an io.Reader implementation.
func (suite *MainTestSuite) TestGetInput() {
	input, err := getInput("stdin", "")
	suite.Nil(err)
	suite.IsType(&os.File{}, input)
	suite.Equal(os.Stdin, input)
	input, err = getInput("file", "main_test.go")
	suite.Nil(err)
	suite.IsType(&os.File{}, input)
	file, _ := input.(*os.File)
	suite.Nil(file.Close())
	input, err = getInput("gzip", "data.gz")
	suite.Nil(err)
	suite.IsType(&gzip.Reader{}, input)
	gzip, _ := input.(*gzip.Reader)
	suite.Nil(gzip.Close())
}

// getInput returns an error if no path is provided for file or gzip input types
func (suite *MainTestSuite) TestGetInputRequiresPathsForFilesAndGZips() {
	var err error
	_, err = getInput("file", "")
	suite.NotNil(err)
	suite.Equal("When input type is 'file', a path must be provided.", err.Error())
	_, err = getInput("gzip", "")
	suite.NotNil(err)
	suite.Equal("When input type is 'gzip', a path must be provided.", err.Error())
}

// getOutput takes an outputType and a path (optionally empty string)
// and returns an io.Writer implementation.
func (suite *MainTestSuite) TestGetOutput() {
	output, err := getOutput("stdout", "")
	suite.Nil(err)
	suite.IsType(&os.File{}, output)
	suite.Equal(os.Stdout, output)
	temp := filepath.Join(os.TempDir(), "test_get_output")
	output, err = getOutput("file", temp)
	suite.Nil(err)
	suite.IsType(&os.File{}, output)
	file, _ := output.(*os.File)
	suite.Nil(file.Close())
	output, err = getOutput("gzip", temp+".gz")
	suite.Nil(err)
	suite.IsType(&gzip.Writer{}, output)
	gzip, _ := output.(*gzip.Writer)
	suite.Nil(gzip.Close())
}

// getOutput returns an error if no path is provided for file or gzip output types
func (suite *MainTestSuite) TestGetOutputRequiresPathsForFilesAndGZips() {
	var err error
	_, err = getOutput("file", "")
	suite.NotNil(err)
	suite.Equal("When output type is 'file', a path must be provided.", err.Error())
	_, err = getOutput("gzip", "")
	suite.NotNil(err)
	suite.Equal("When output type is 'gzip', a path must be provided.", err.Error())
}

// NaturalSort reads from an io.Reader, sorts the input naturally and
// writes to an io.Writer.
func (suite *MainTestSuite) TestNaturalSort() {
	input := bytes.NewBufferString("02z a2 a01 a010 a10z 1z")
	output := bytes.NewBuffer([]byte{})
	err := NaturalSort(input, output)
	suite.Nil(err)
	suite.Equal("1z\n02z\na01\na2\na010\na10z\n", output.String())
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
