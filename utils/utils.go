//Some additional utilities that are useful when processing CSV headers and data.
package utils

import (
	"math/rand"
	"time"
)

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

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}
