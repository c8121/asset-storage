package util

import "strconv"

func Atoi(s string, returnOnFail int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return returnOnFail
	}
	return i
}
