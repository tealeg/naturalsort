package main

// ByNaturalOrder implements the sort.Interface for []item.
type ByNaturalOrder []item

func (items ByNaturalOrder) Len() int {
	return len(items)
}

func (items ByNaturalOrder) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items ByNaturalOrder) Less(i, j int) bool {
	itemA := items[i]
	itemB := items[j]
	lenA := len(itemA.parts)
	lenB := len(itemB.parts)
	for index := 0; index < lenA && index < lenB; index++ {
		partA := itemA.parts[index]
		partB := itemB.parts[index]
		if partA.LessThan(partB) {
			return true
		}
		if partB.LessThan(partA) {
			return false
		}
	}
	return lenA < lenB
}
