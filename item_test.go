package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ItemTestSuite struct {
	suite.Suite
}

// Attempting to Read to a zero length buffer should return 0 bytes read
func (suite *ItemTestSuite) TestReadToZeroLengthBuffer() {
	buff := make([]byte, 0, 0)
	item := item{parts: partList{runePart{runeVal: 'a'}}}
	count, err := item.Read(buff)
	suite.Equal(0, count)
	suite.Nil(err)
}

// Attempting to Read more bytes than are present in the item returns an io.EOF
func (suite *ItemTestSuite) TestReadBeyondLengthOfItemReturnsEOF() {
	buff := make([]byte, 2, 2)
	itm := item{parts: partList{runePart{runeVal: 'a'}}}
	count, err := itm.Read(buff)
	suite.Equal(1, count)
	suite.NotNil(err)
	suite.Equal(io.EOF, err)
	suite.Equal('a', rune(buff[0]))
}

// Partial read of numeric sequence is valid
func (suite *ItemTestSuite) TestPartialRead() {
	buff := make([]byte, 1, 1)
	part, _ := newIntPartFromString("012")
	itm := item{parts: partList{part}}
	count, err := itm.Read(buff)
	suite.Nil(err)
	suite.Equal(1, count)
	suite.Equal('0', rune(buff[0]))
}

// Read will read all of the items parts in order
func (suite *ItemTestSuite) TestReadAllSubParts() {
	buff := make([]byte, 5, 5)
	part0 := runePart{runeVal: 'a'}
	part1, _ := newIntPartFromString("010")
	part2 := runePart{runeVal: 'b'}
	item := item{parts: partList{part0, part1, part2}}
	count, err := item.Read(buff)
	suite.Equal(5, count)
	suite.Nil(err)
	suite.Equal("a010b", string(buff))
}

func TestItemTestSuite(t *testing.T) {
	suite.Run(t, new(ItemTestSuite))
}
