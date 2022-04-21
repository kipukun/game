package object

import "sort"

var origin *obj

func init() {
	origin = new(obj)
}

type Collection []Object

func (c Collection) Len() int {
	return len(c)
}

func NewCollection(os ...Object) Collection {
	return Collection(os)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (c Collection) Less(i int, j int) bool {
	return Distance(origin, c[i]) < Distance(origin, c[j])
}

// Swap swaps the elements with indexes i and j.
func (c Collection) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Collection) Sort() {
	sort.Stable(c)
}
