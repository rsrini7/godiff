package main

//
// An O(ND) Difference Algorithm: Find middle snake
//
func algorithm_sms(data1, data2 []int, v []int) (int, int, int, int) {

	end1, end2 := len(data1), len(data2)
	max := end1 + end2 + 1
	up_k := end1 - end2
	odd := (up_k & 1) != 0
	down_off, up_off := max, max-up_k+max+max+2

	v[down_off+1] = 0
	v[down_off] = 0
	v[up_off+up_k-1] = end1
	v[up_off+up_k] = end1

	var k, x, u, z int

	for d := 1; true; d++ {
		up_k_plus_d := up_k + d
		up_k_minus_d := up_k - d
		for k = -d; k <= d; k += 2 {
			x = v[down_off+k+1]
			if k > -d && (k == d || z >= x) {
				x, z = z+1, x
			} else {
				z = x
			}
			for u = x; x < end1 && x-k < end2 && data1[x] == data2[x-k]; x++ {
			}
			if odd && (up_k_minus_d < k) && (k < up_k_plus_d) && v[up_off+k] <= x {
				return u, u - k, x, x - k
			}
			v[down_off+k] = x
		}
		z = v[up_off+up_k_minus_d-1]
		for k = up_k_minus_d; k <= up_k_plus_d; k += 2 {
			x = z
			if k < up_k_plus_d {
				z = v[up_off+k+1]
				if k == up_k_minus_d || z <= x {
					x = z - 1
				}
			}
			for u = x; x > 0 && x > k && data1[x-1] == data2[x-k-1]; x-- {
			}
			if !odd && (-d <= k) && (k <= d) && x <= v[down_off+k] {
				return x, x - k, u, u - k
			}
			v[up_off+k] = x
		}
	}
	return 0, 0, 0, 0 // should not reach here
}

//
// Special case for algorithm_sms() with only 1 item.
//
func find_one_sms(value int, list []int) (int, int) {
	for i, v := range list {
		if v == value {
			return 0, i
		}
	}
	return 1, 0
}

//
// An O(ND) Difference Algorithm: Find LCS
//
func algorithm_lcs(data1, data2 []int, change1, change2 []bool, v []int) {

	start1, start2 := 0, 0
	end1, end2 := len(data1), len(data2)

	// matches found at start and end of list
	for start1 < end1 && start2 < end2 && data1[start1] == data2[start2] {
		start1++
		start2++
	}
	for start1 < end1 && start2 < end2 && data1[end1-1] == data2[end2-1] {
		end1--
		end2--
	}

	len1, len2 := end1-start1, end2-start2

	switch {
	case len1 == 0:
		for start2 < end2 {
			change2[start2] = true
			start2++
		}

	case len2 == 0:
		for start1 < end1 {
			change1[start1] = true
			start1++
		}

	case len1 == 1 && len2 == 1:
		change1[start1] = true
		change2[start2] = true

	default:
		data1, change1 = data1[start1:end1], change1[start1:end1]
		data2, change2 = data2[start2:end2], change2[start2:end2]

		var x0, y0, x1, y1 int

		if len(data1) == 1 {
			// match one item, use simple search function
			x0, y0 = find_one_sms(data1[0], data2)
			x1, y1 = x0, y0
		} else if len(data2) == 1 {
			// match one item, use simple search function
			y0, x0 = find_one_sms(data2[0], data1)
			x1, y1 = x0, y0
		} else {
			// Find a point with the longest common sequence
			x0, y0, x1, y1 = algorithm_sms(data1, data2, v)
		}

		// Use the partitions to split this problem into subproblems.
		algorithm_lcs(data1[:x0], data2[:y0], change1[:x0], change2[:y0], v)
		algorithm_lcs(data1[x1:], data2[y1:], change1[x1:], change2[y1:], v)
	}
}
