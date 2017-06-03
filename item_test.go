package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ItemTestSuite struct {
	suite.Suite
}

// Attempting to Read to a zero length buffer should return 0 bytes read and io.EOF
func (suite *ItemTestSuite) TestReadToZeroLengthBuffer() {
	buff := make([]byte, 0, 0)
	item := item{parts: partList{runePart('a')}}
	count, err := item.Read(buff)
	suite.Equal(0, count)
	suite.NotNil(err)
	suite.Equal(io.EOF, err)
}

func TestItemTestSuite(t *testing.T) {
	suite.Run(t, new(ItemTestSuite))
}
