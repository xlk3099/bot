package utils

import "strconv"

// Float64ToString converts a float64 number to string.
func Float64ToString(input float64) string {
	return strconv.FormatFloat(input, 'f', 6, 64)
}
