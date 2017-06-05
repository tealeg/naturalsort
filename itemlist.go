package main

import (
	"bytes"
	"io"
	"sort"
)

// itemList wraps an array of items and supports the
// io.ReadWriter interface.
type itemList struct {
	items     []item
	nextItem  item
	numBuff   bytes.Buffer
	readIndex int
}

// Pull the string representing a number off the numBuff and convert
// it to an intPart on itemList.parts
func (il *itemList) addIntPart() (err error) {
	if il.numBuff.Len() > 0 {
		var iPart part
		iPart, err = newIntPartFromString(il.numBuff.String())
		if err != nil {
			return
		}
		il.nextItem.parts = append(il.nextItem.parts, iPart)
		il.numBuff.Reset()
	}
	return
}

// Implementation of io.Writer for itemList
func (il *itemList) Write(b []byte) (n int, err error) {

	n = 0

	for _, char := range b {
		switch {
		case char < 33:
			// Whitespace, newline or any control character is treated as a terminator
			err = il.addIntPart()
			if err != nil {
				return
			}

			// We'll 'normalise' all separators to newline for our purposes.
			il.nextItem.parts = append(il.nextItem.parts, runePart{runeVal: '\n'})
			il.items = append(il.items, il.nextItem)
			il.nextItem = item{}
		case char > 47 && char < 58:
			// This is a rune representing a number
			il.numBuff.WriteByte(char)
		default:
			err = il.addIntPart()
			if err != nil {
				return
			}
			il.nextItem.parts = append(il.nextItem.parts, runePart{runeVal: rune(char)})
		}
		n++
	}
	return
}

// Flush unterminated input to the itemList
func (il *itemList) Flush() error {
	if il.numBuff.Len() > 0 {
		// Flush any outstanding input to the items array
		err := il.addIntPart()
		if err != nil {
			return err
		}
	}
	if len(il.nextItem.parts) > 0 {
		il.nextItem.parts = append(il.nextItem.parts, runePart{runeVal: '\n'})
		il.items = append(il.items, il.nextItem)
	}
	return nil
}

// Implementation of io.Reader for itemList
func (il *itemList) Read(b []byte) (n int, err error) {
	const outputSeparator = '\n'
	count := 0
	buffLen := len(b)
	for {
		if n == buffLen {
			return
		}
		itm := il.items[il.readIndex]
		count, err = itm.Read(b[n:])
		n += count
		if err == io.EOF {
			// We've exhausted this item and we should
			// look at the next, but, if we don't hit EOF
			// here it means we've run out of space in the
			// buffer we're reading to, so we'll want to
			// carry on from here (which is what happens
			// when we don't increment the readIndex).
			il.readIndex++
			// err is defined in the signature scope, so
			// we'd better clear it so we don't
			// accidentally return io.EOF at the end of
			// the function.
			if il.readIndex < len(il.items) {
				err = nil
				continue
			}
			break
		}
		if err != nil {
			break
		}
	}
	return
}

// Sort the itemList.items
func (il *itemList) Sort() {
	sort.Sort(ByNaturalOrder(il.items))
}
