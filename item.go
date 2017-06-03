package main

import "io"

type partList []part

type item struct {
	parts partList
	index int
}

// Implement io.Reader for item
func (i *item) Read(p []byte) (n int, err error) {

	if len(p) == 0 {
		return 0, io.EOF
	}
	return
	// r.prevRune = -1
	// n = copy(b, r.s[r.i:])
	// r.i += int64(n)
	// return
	// if len(p) == 0 {
	// 	n = 0
	// 	err = fmt.Errorf("Cannot read item into zero length buffer.")
	// 	return
	// }
	// p[
}
