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

// Termination in successive calls to Write causes carried over items
// to be pushed onto the itemList.
func (suite *ItemListTestSuite) TestMultipleWrites() {
	il := itemList{}
	input1 := []byte{'a', 'b'}

	count, err := il.Write(input1)
	suite.Nil(err)
	suite.Equal(2, count)
	suite.Equal(0, len(il.items))

	input2 := []byte{' '}

	count, err = il.Write(input2)
	suite.Nil(err)
	suite.Equal(1, count)
	suite.Equal(1, len(il.items))

	// Check the value of the item pushed to the list
	buff := make([]byte, 2, 2)

	count, err = il.items[0].Read(buff)
	suite.Nil(err)
	suite.Equal(2, count)
	suite.Equal("ab", string(buff))
}

// Closing the itemList causes the last unterminated item to be pushed to the itemList
func (suite *ItemListTestSuite) TestCloseFlushesNextItem() {
	il := itemList{}
	input := []byte{'a', 'b'}

	count, err := il.Write(input)
	suite.Nil(err)
	suite.Equal(2, count)
	suite.Equal(0, len(il.items))

	err = il.Close()
	suite.Nil(err)

	suite.Equal(1, len(il.items))

	// Check the value of the item pushed to the list
	buff := make([]byte, 2, 2)

	count, err = il.items[0].Read(buff)
	suite.Nil(err)
	suite.Equal(2, count)
	suite.Equal("ab", string(buff))
}

// Closing the itemList with no outstanding unterminated items doesn't add an empty item to the itemList
func (suite *ItemListTestSuite) TestCloseWithoutUnterminatedItemsIsANoOp() {
	il := itemList{}
	err := il.Close()
	suite.Nil(err)
	suite.Equal(0, len(il.items))
}

func TestItemListTestSuite(t *testing.T) {
	suite.Run(t, new(ItemListTestSuite))
}
