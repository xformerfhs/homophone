package randomlist

import (
	"math/rand/v2"
)

// ******** Public types ********

// RandomList implements a list of elements that are accessed in a random order.
type RandomList[T any] struct {
	baseSlice []T
	length    int
	index     []int
	actIndex  int
}

// ******** Public generation function ********

// New creates a new random list.
func New[T any](s []T) *RandomList[T] {
	sliceLen := len(s)
	index := make([]int, sliceLen)
	if sliceLen > 1 {
		newRandomIndexList(index, sliceLen)
	}

	return &RandomList[T]{
		baseSlice: s,
		length:    sliceLen,
		index:     index,
		actIndex:  0,
	}
}

// ******** Public functions ********

// Len returns the length of the random list.
func (r *RandomList[T]) Len() int {
	return r.length
}

// BaseList returns the base list of the random list.
func (r *RandomList[T]) BaseList() []T {
	return r.baseSlice
}

// RandomElement returns a random element from the list.
func (r *RandomList[T]) RandomElement() T {
	sliceLen := r.length
	baseSlice := r.baseSlice
	// Special handling if the slice has at most one element.
	// This will panic if the base slice is empty.
	if sliceLen <= 1 {
		return baseSlice[0]
	}

	var resultIndex int
	resultIndex, r.actIndex = incIndex(r.index, r.actIndex, sliceLen)

	// Return the random element.
	return baseSlice[resultIndex]
}

// ******** Private functions ********

// incIndex returns the current random index and increments the index into the index slice.
func incIndex(index []int, actIndex int, sliceLen int) (int, int) {
	if actIndex >= sliceLen {
		newRandomIndexList(index, sliceLen)
		actIndex = 0
	}

	result := index[actIndex]

	actIndex++

	return result, actIndex
}

// newRandomIndexList fills the index slice with a new random shuffle of the indices.
func newRandomIndexList(index []int, count int) {
	_ = index[count-1] // Skip index check in loop

	for i := 0; i < count; i++ {
		index[i] = i
	}

	rand.Shuffle(count, func(i, j int) { index[i], index[j] = index[j], index[i] })
}
