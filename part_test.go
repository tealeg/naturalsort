package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PartTestSuite struct {
	suite.Suite
}

// LessThan handles alphabetical chars
func (suite *PartTestSuite) TestLessThanAlpha() {
	part1 := runePart('z')
	part2 := runePart('x')
	part3 := runePart('x')
	suite.False(part1.LessThan(part2))
	suite.True(part2.LessThan(part1))
	suite.False(part2.LessThan(part3))
	suite.False(part3.LessThan(part2))
}

func TestPartTestSuite(t *testing.T) {
	suite.Run(t, new(PartTestSuite))
}
