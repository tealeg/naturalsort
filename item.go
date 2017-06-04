package main

import (
	"io"
	"strings"
)

// Convenience type
type partList []part

type item struct {
	parts partList
	index int
}

// Implement io.Reader for item
func (i *item) Read(b []byte) (n int, err error) {
	count := 0
	index := 0
	buffLen := len(b)
	for _, part := range i.parts[i.index:] {
		if index >= buffLen {
			return
		}
		count, err = part.Read(b[index:])
		if err != nil && err != io.EOF {
			i.index += n
			return
		}
		n += count
		index += count
	}
	i.index += n
	return
}

func (i *item) String() string {
	s := []string{}
	for _, part := range i.parts {
		s = append(s, part.String())
	}
	return strings.Join(s, "")

}
