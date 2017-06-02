package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NaturalSortTestSuite struct {
	suite.Suite
}

// When no numeric elements are present, alphabetical sorting is used.
func (suite *NaturalSortTestSuite) TestSortAlpha() {
	item0 := item{runePart('e'), runePart('d'), runePart('f')}
	item1 := item{runePart('a'), runePart('b'), runePart('z')}
	item2 := item{runePart('a'), runePart('e'), runePart('g')}
	items := []item{item0, item1, item2}
	sort.Sort(ByNaturalOrder(items))
	suite.Equal(item1, items[0])
	suite.Equal(item2, items[1])
	suite.Equal(item0, items[2])
}

func TestNaturalSortTestSuite(t *testing.T) {
	suite.Run(t, new(NaturalSortTestSuite))
}
