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
	buffLen := len(b)
	for {
		if n >= buffLen {
			break
		}
		part := i.parts[i.index]
		count, err = part.Read(b[n:])
		n += count
		if err == io.EOF {
			i.index++
			if i.index < len(i.parts) {
				err = nil
				continue
			}
			break
			// if there are no more parts, we'll exit just below with an EOF.
		}
		if err != nil {
			break
		}
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
