package main

import (
	"io"
	"strconv"
	"strings"
)

type part interface {
	io.Reader
	LessThan(other part) bool
	String() string
}

type runePart struct {
	runeVal rune
	read    bool
}

func (r runePart) LessThan(other part) bool {
	if otherR, ok := other.(runePart); ok {
		return r.runeVal < otherR.runeVal
	}
	// The other part must be an intPart, which will always come first
	return false
}

func (r runePart) String() string {
	return string(r.runeVal)
}

func (r runePart) Read(b []byte) (n int, err error) {

	buffLen := len(b)
	if buffLen == 0 {
		return
	}
	if r.read {
		err = io.EOF
		return
	}
	b[0] = byte(r.runeVal)
	n = 1
	r.read = true
	if buffLen > 1 {
		err = io.EOF
	}
	return
}

type intPart struct {
	intVal int
	strVal string
	reader *strings.Reader
}

func newIntPartFromString(s string) (i intPart, err error) {
	i = intPart{strVal: s}
	i.intVal, err = strconv.Atoi(s)
	if err == nil {
		i.reader = strings.NewReader(i.strVal)
	}
	return
}

func (i intPart) String() string {
	if len(i.strVal) == 0 {
		return strconv.Itoa(i.intVal)
	}
	return i.strVal
}

func (i intPart) LessThan(other part) bool {
	if otherI, ok := other.(intPart); ok {
		return i.intVal < otherI.intVal
	}
	// The other part must be a runePart, which will always come second
	return true
}

func (i intPart) Read(b []byte) (n int, err error) {
	return i.reader.Read(b)
}
