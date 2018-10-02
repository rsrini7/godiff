package main

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
func to_lower_byte(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b - 'A' + 'a'
	}
	return b
}

//
// Test for space character
//
func is_space(b byte) bool {
	return b == ' ' || b == '\t' || b == '\v' || b == '\f'
}

//RandShuffle to shuffle the given slices of numbers
func RandShuffle(input []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(input), func(i, j int) { input[i], input[j] = input[j], input[i] })
}
