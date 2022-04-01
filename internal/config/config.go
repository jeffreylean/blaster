package config

import (
	"fmt"
	"os"
	"strconv"
)

func GetConfigInt64(in string) int64 {
	n := os.Getenv(in)
	var out int64
	if n != "" {
		i, err := strconv.Atoi(n)
		if err != nil {
			fmt.Println(err)
		}
		out = int64(i)
	}
	return out
}

func GetConfigString(in string) string {
	out := os.Getenv(in)
	return out
}
