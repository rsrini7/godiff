package main

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
