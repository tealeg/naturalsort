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

	// Check that the nextItem contains all the input
	buff := make([]byte, 3, 3)
	count, err = il.nextItem.Read(buff)
	suite.Nil(err)
	suite.Equal(3, count)
	suite.Equal("abc", string(buff))
}

// Byte slices that include a terminator push new items to the itemList
func (suite *ItemListTestSuite) TestWriteTerminatedInput() {
	il := itemList{}
	input := []byte{'a', 'b', ' ', 'c', 'd', '\n', 'e', 'f'}
	count, err := il.Write(input)
	suite.Nil(err)
	suite.Equal(8, count)
	suite.Equal(2, len(il.items))

	// Check the content of the items
	buff := make([]byte, 2, 2)

	count, err = il.items[0].Read(buff)
	suite.Nil(err)
	suite.Equal(2, count)
	suite.Equal("ab", string(buff))

	count, err = il.items[1].Read(buff)
	suite.Nil(err)
	suite.Equal(2, count)
	suite.Equal("cd", string(buff))

	// Also check the nextItem
	count, err = il.nextItem.Read(buff)
	suite.Nil(err)
	suite.Equal(2, count)
	suite.Equal("ef", string(buff))
}

func TestItemListTestSuite(t *testing.T) {
	suite.Run(t, new(ItemListTestSuite))
}
