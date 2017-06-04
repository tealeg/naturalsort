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
func (i *item) Read(p []byte) (n int, err error) {
	count := 0
	index := 0
	buffLen := len(p)
	for _, part := range i.parts {
		if index >= buffLen {
			return
		}
		count, err = part.Read(p[index:])
		if err != nil && err != io.EOF {
			return
		}
		n += count
		index += count
	}
	return
}

func (i *item) String() string {
	s := []string{}
	for _, part := range i.parts {
		s = append(s, part.String())
	}
	return strings.Join(s, "")

}
