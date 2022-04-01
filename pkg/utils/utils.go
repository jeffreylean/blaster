package utils

import "strconv"

func ToInt(in string) int {
	int, err := strconv.Atoi(in)
	if err != nil {
		return 0
	}
	return int
}
