package main

type part interface {
	LessThan(other part) bool
}

type runePart rune

func (r runePart) LessThan(other part) bool {
	if otherR, ok := other.(runePart); ok {
		return r < otherR
	}
	// The other part must be an intPart, which will always come first
	return false
}

type intPart int

func (r intPart) LessThan(other part) bool {
	if otherR, ok := other.(intPart); ok {
		return r < otherR
	}
	// The other part must be a runePart, which will always come second
	return true
}
