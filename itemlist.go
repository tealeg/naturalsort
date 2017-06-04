package main

import (
	"bytes"
	"io"
)

// itemList wraps an array of items and supports the
// io.ReadWriteCloser interface.
type itemList struct {
	items     []item
	nextItem  item
	numBuff   bytes.Buffer
	readIndex int
}

// Implementation of io.Writer for itemList
func (il *itemList) Write(b []byte) (n int, err error) {
	var iPart part
	n = 0
	for _, char := range b {
		switch {
		case char < 33:
			// Whitespace, newline or any control character is treated as a terminator
			il.items = append(il.items, il.nextItem)
			il.nextItem = item{}
		case char > 47 && char < 58:
			// This is a rune representing a number
			il.numBuff.WriteByte(char)
		default:
			if il.numBuff.Len() > 0 {
				iPart, err = newIntPartFromString(il.numBuff.String())
				if err != nil {
					return
				}
				il.nextItem.parts = append(il.nextItem.parts, iPart)
				il.numBuff.Reset()
			}
			il.nextItem.parts = append(il.nextItem.parts, runePart{runeVal: rune(char)})
		}
		n++
	}
	return
}

// Implementation of io.Closer for itemList
func (il *itemList) Close() error {
	if len(il.nextItem.parts) > 0 {
		// Flush any outstanding input to the items array
		il.items = append(il.items, il.nextItem)
	}
	return nil
}

// Implementation of io.Reader for itemList
func (il *itemList) Read(b []byte) (n int, err error) {
	const outputSeparator = '\n'
	count := 0
	index := 0
	buffLen := len(b)
	writeableLen := buffLen - 1 // We need char to write our output separator to!
	for _, itm := range il.items[il.readIndex:] {
		if index >= writeableLen {
			return
		}
		count, err = itm.Read(b[index:])
		if err == io.EOF {
			// We've exhausted this item and we should
			// look at the next, but, if we don't hit EOF
			// here it means we've run out of space in the
			// buffer we're reading to, so we'll want to
			// carry on from here (which is what happens
			// when we don't increment the readIndex).
			il.readIndex++
			// We have to tack separators back into the
			// output, because we threw away the ones from
			// the input.
			b[index+count] = byte(outputSeparator)
			// .. which means we have bump the count as well.
			count++
			// err is defined in the signature scope, so
			// we'd better clear it so we don't
			// accidentally return io.EOF at the end of
			// the function.
			err = nil
		}
		if err != nil {
			return
		}
		n += count
		index += count
	}
	return
}
