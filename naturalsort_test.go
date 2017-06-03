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

// When no alphabetical runes are present, numerical sorting is used
func (suite *NaturalSortTestSuite) TestSortNumeric() {
	item0 := item{intPart(700)}
	item1 := item{intPart(80)}
	item2 := item{intPart(9)}
	items := []item{item0, item1, item2}
	sort.Sort(ByNaturalOrder(items))
	suite.Equal(item2, items[0])
	suite.Equal(item1, items[1])
	suite.Equal(item0, items[2])
}

// When both runeParts and intParts are present, natural ordering is achieved.
func (suite *NaturalSortTestSuite) TestSortMixed() {
	item0 := item{intPart(700)}
	item1 := item{runePart('a')}
	item2 := item{intPart(80)}
	item3 := item{runePart('b')}
	items := []item{item0, item1, item2, item3}
	sort.Sort(ByNaturalOrder(items))
	suite.Equal(item2, items[0])
	suite.Equal(item0, items[1])
	suite.Equal(item1, items[2])
	suite.Equal(item3, items[3])
}

// When an item is an exact prefix os the item it is compared to then
// it is sorted ahead of the other item.
func (suite *NaturalSortTestSuite) TestSortSubItems() {
	item0 := item{runePart('a'), intPart(1)}
	item1 := item{runePart('a')}
	item2 := item{runePart('a'), intPart(1), runePart('a')}
	items := []item{item0, item1, item2}
	sort.Sort(ByNaturalOrder(items))
	suite.Equal(item1, items[0])
	suite.Equal(item0, items[1])
	suite.Equal(item2, items[2])
}

// When items in the item list are compound of intParts and runeParts,
// natural ordering is achieved.
func (suite *NaturalSortTestSuite) TestSortCompounds() {
	item0 := item{runePart('a'), intPart(1)}
	item1 := item{runePart('a'), intPart(0), runePart('a')}
	item2 := item{intPart(2), runePart('a')}
	item3 := item{runePart('b'), intPart(9)}
	item4 := item{intPart(10), runePart('a')}
	item5 := item{runePart('b'), intPart(80)}
	item6 := item{runePart('a'), intPart(0), runePart('a'), runePart('b')}
	items := []item{item0, item1, item2, item3, item4, item5, item6}
	sort.Sort(ByNaturalOrder(items))
	suite.Equal(item2, items[0])
	suite.Equal(item4, items[1])
	suite.Equal(item1, items[2])
	suite.Equal(item6, items[3])
	suite.Equal(item0, items[4])
	suite.Equal(item3, items[5])
	suite.Equal(item5, items[6])
}

func TestNaturalSortTestSuite(t *testing.T) {
	suite.Run(t, new(NaturalSortTestSuite))
}
