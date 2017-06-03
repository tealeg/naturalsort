package main

type part interface {
	LessThan(other part) bool
}

type runePart rune

func (r runePart) LessThan(other part) bool {
	if otherR, ok := other.(runePart); ok {
		return r < otherR
	}
	return false
}

type intPart int

func (r intPart) LessThan(other part) bool {
	if otherR, ok := other.(intPart); ok {
		return r < otherR
	}
	return false
}
