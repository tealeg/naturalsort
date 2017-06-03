package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RunePartTestSuite struct {
	suite.Suite
}

// LessThan handles rune vs rune comparisons
func (suite *RunePartTestSuite) TestLessThanRuneVRune() {
	part1 := runePart('z')
	part2 := runePart('x')
	part3 := runePart('x')
	suite.False(part1.LessThan(part2))
	suite.True(part2.LessThan(part1))
	suite.False(part2.LessThan(part3))
	suite.False(part3.LessThan(part2))
}

func TestRunePartTestSuite(t *testing.T) {
	suite.Run(t, new(RunePartTestSuite))
}

type IntPartTestSuite struct {
	suite.Suite
}

// LessThan handles int vs int comparisons
func (suite *IntPartTestSuite) TestLessThanIntVInt() {
	part0 := intPart(9)
	part1 := intPart(80)
	part2 := intPart(700)
	suite.True(part0.LessThan(part1))
	suite.True(part1.LessThan(part2))
	suite.False(part0.LessThan(part0))
	suite.False(part1.LessThan(part1))

}

func TestIntPartTestSuite(t *testing.T) {
	suite.Run(t, new(IntPartTestSuite))
}
