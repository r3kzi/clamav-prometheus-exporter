package clamav

import (
	"fmt"
	"strconv"
)

func toFloat(s string) float64 {
	float, err := strconv.ParseFloat(s, 64)
	if err != nil {
		_ = fmt.Errorf("couldn't parse string to float: %s", err)
	}
	return float
}
