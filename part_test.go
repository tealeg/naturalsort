package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RunePartTestSuite struct {
	suite.Suite
}

// LessThan handles rune vs rune comparisons
func (suite *RunePartTestSuite) TestLessThanRuneVsRune() {
	part1 := runePart('z')
	part2 := runePart('x')
	part3 := runePart('x')
	suite.False(part1.LessThan(part2))
	suite.True(part2.LessThan(part1))
	suite.False(part2.LessThan(part3))
	suite.False(part3.LessThan(part2))
}

// LessThan handles rune vs int comparisons, ints are always less than runes
func (suite *RunePartTestSuite) TestLessThanRuneVsInt() {
	part0 := intPart{intVal: 99}
	part1 := runePart('a')
	part2 := intPart{intVal: 99}
	suite.False(part1.LessThan(part2))
	suite.True(part0.LessThan(part1))
}

// Calling runePart.Read with a zero length buffer returns a 0 count
func (suite *RunePartTestSuite) TestReadToZeroLengthBuffer() {
	part := runePart('a')
	buff := make([]byte, 0, 0)
	count, err := part.Read(buff)
	suite.Nil(err)
	suite.Equal(0, count)
	suite.Equal("", string(buff))
}

// Calling runePart.Read with a 1 byte buffer returns a 1 count
func (suite *RunePartTestSuite) TestReadTo1ByteBuffer() {
	part := runePart('a')
	buff := make([]byte, 1, 1)
	count, err := part.Read(buff)
	suite.Nil(err)
	suite.Equal(1, count)
	suite.Equal("a", string(buff))
}

// Calling runePart.Read with a multi byte buffer returns a 1 count, and an io.EOF error
func (suite *RunePartTestSuite) TestReadToMultiByteBuffer() {
	part := runePart('a')
	buff := make([]byte, 2, 2)
	count, err := part.Read(buff)
	suite.NotNil(err)
	suite.Equal(io.EOF, err)
	suite.Equal(1, count)
	suite.Equal('a', rune(buff[0]))
}

func TestRunePartTestSuite(t *testing.T) {
	suite.Run(t, new(RunePartTestSuite))
}

type IntPartTestSuite struct {
	suite.Suite
}

// newIntPartFromString returns an intPart representing the integer
// provided as a string.
func (suite *IntPartTestSuite) TestNewIntPartFromString() {
	part, err := newIntPartFromString("100")
	suite.Nil(err)
	suite.Equal(100, part.intVal)
	suite.Equal("100", part.strVal)
	suite.NotNil(part.reader)
}

// newIntPartFromString return an error if it is passed a value that
// can be converted to an integer.
func (suite *IntPartTestSuite) TestNewIntPartFromStringReturnsErrorsOnBadInput() {
	_, err := newIntPartFromString("shoe")
	suite.NotNil(err)
}

// LessThan handles int vs int comparisons
func (suite *IntPartTestSuite) TestLessThanIntVsInt() {
	part0 := intPart{intVal: 9}
	part1 := intPart{intVal: 80}
	part2 := intPart{intVal: 700}
	suite.True(part0.LessThan(part1))
	suite.True(part1.LessThan(part2))
	suite.False(part0.LessThan(part0))
	suite.False(part1.LessThan(part1))

}

// LessThan handles int vs rune comparisons, ints are always less than runes
func (suite *IntPartTestSuite) TestLessThanIntVsRune() {
	part0 := runePart('a')
	part1 := intPart{intVal: 99}
	part2 := runePart('a')
	suite.True(part1.LessThan(part2))
	suite.False(part0.LessThan(part1))
}

// Calling intPart.Read with a zero length buffer causes 0 bytes to be read
func (suite *IntPartTestSuite) TestReadToZeroLengthBuffer() {
	part, _ := newIntPartFromString("9")
	buff := make([]byte, 0, 0)
	count, _ := part.Read(buff)
	suite.Equal(0, count)
}

// Calling intPart.Read will read the strVal of the intPart
func (suite *IntPartTestSuite) TestRead() {
	part, err := newIntPartFromString("001")
	suite.Nil(err)
	buff := make([]byte, 3, 3)
	count, err := part.Read(buff)
	suite.Nil(err)
	suite.Equal(3, count)
	suite.Equal("001", string(buff))
}

func TestIntPartTestSuite(t *testing.T) {
	suite.Run(t, new(IntPartTestSuite))
}
