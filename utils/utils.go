//Some additional utilities that are useful when processing CSV headers and data.
package utils

import (
	"math/rand"
	"time"
)

type Index map[string]int

// Return a map that maps each string in the input slice to its index in the slice.
func NewIndex(a []string) Index {
	index := make(map[string]int)
	for i, v := range a {
		index[v] = i
	}
	return index
}

// Answer true if the index contains the specified string.
func (i Index) Contains(k string) bool {
	_, ok := i[k]
	return ok
}

// Calculate the intersection between two string slices. The first returned slice
// is the intersection between the two slices. The second returned slice is
// a slice of elements in the first slice but not the second. The third returned
// slice is a slice of elements in the second slice but not the first.
func Intersect(a []string, b []string) ([]string, []string, []string) {
	index := NewIndex(a)
	result := make([]string, 0, len(b))
	aNotB := make([]string, len(a), len(a))
	copy(aNotB, a)
	bNotA := make([]string, 0, len(b))
	for _, v := range b {
		if i, ok := index[v]; ok {
			result = append(result, v)
			aNotB[i] = ""
		} else {
			bNotA = append(bNotA, v)
		}
	}
	i := 0
	for j := range a {
		present := (aNotB[j] == a[j])
		aNotB[i] = a[j]
		if present {
			i++
		}
	}
	aNotB = aNotB[0:i]
	return result, aNotB, bNotA
}

// shortcut functions. hopefully will be inlined by compiler
func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// shortcut functions. hopefully will be inlined by compiler
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// convert byte to lower case
func ToLowerByte(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b - 'A' + 'a'
	}
	return b
}

//
// Test for space character
//
func IsSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\v' || b == '\f'
}

//RandShuffle to shuffle the given slices of numbers
func RandShuffle(input []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(input), func(i, j int) { input[i], input[j] = input[j], input[i] })
}

//Equal :test whether the given two string slices are equal
func Equal(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Find returns the smallest index i at which x == a[i],
// or -1 if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}
