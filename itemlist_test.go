package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ItemListTestSuite struct {
	suite.Suite
}

// Byte slices that don't include a terminator
func (suite *ItemListTestSuite) TestWriteUnterminatedInputToBuffer() {
	il := itemList{}
	input := []byte{'a', 'b', 'c'}
	count, err := il.Write(input)
	suite.Nil(err)
	suite.Equal(3, count)
	// No items are complete
	suite.Equal(0, len(il.items))
	suite.Equal(3, len(il.nextItem.parts))
}

func TestItemListTestSuite(t *testing.T) {
	suite.Run(t, new(ItemListTestSuite))
}
