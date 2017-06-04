package main

import (
	"bytes"
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
	suite.Equal("abc", il.nextItem.String())
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
	suite.Equal("ab\n", il.items[0].String())
	suite.Equal("cd\n", il.items[1].String())
	// nextItem is not terminated and therefore has no newline
	suite.Equal("ef", il.nextItem.String())
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
	suite.Equal("ab\n", il.items[0].String())
}

// Numeric input is bunched into intParts
func (suite *ItemListTestSuite) TestNumericInputIsGrouped() {
	il := itemList{}
	input := []byte{'a', '0', '1', '2', 'b'}
	count, err := il.Write(input)
	suite.Nil(err)
	suite.Equal(5, count)
	err = il.Close()
	suite.Nil(err)
	suite.Equal(1, len(il.items))
	item := il.items[0]
	// There are four parts because Close() adds a separator (newline)
	suite.Equal(4, len(item.parts))
	suite.Equal("a012b\n", item.String())

}

// Closing the itemList causes the last unterminated item to be pushed
// to the itemList (and a separator is added).
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
	suite.Equal("ab\n", il.items[0].String())
}

// Closing the itemList with no outstanding unterminated items doesn't add an empty item to the itemList
func (suite *ItemListTestSuite) TestCloseWithoutUnterminatedItemsIsANoOp() {
	il := itemList{}
	err := il.Close()
	suite.Nil(err)
	suite.Equal(0, len(il.items))
}

// Reading the itemList to exhaustion causes all it's items to be read, sequentially.
func (suite *ItemListTestSuite) TestReadSequentiallyReadsItems() {
	input := bytes.NewBufferString("abc123 abc234 123xyz")
	il := itemList{}
	count, err := il.Write(input.Bytes())
	suite.Nil(err)
	suite.Equal(20, count)
	err = il.Close()
	suite.Equal(3, len(il.items))
	output := make([]byte, 21, 21)
	count, err = il.Read(output)
	suite.Nil(err)
	// We have one more character in the output because we've
	// forced a newline on the end.
	suite.Equal(21, count)
	suite.Equal("abc123\nabc234\n123xyz\n", string(output))
}

func TestItemListTestSuite(t *testing.T) {
	suite.Run(t, new(ItemListTestSuite))
}
