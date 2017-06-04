package main

import "bytes"

// itemList wraps an array of items and supports the
// io.ReadWriteCloser interface.
type itemList struct {
	items    []item
	nextItem item
	numBuff  bytes.Buffer
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
		il.items = append(il.items, il.nextItem)
	}
	return nil
}
