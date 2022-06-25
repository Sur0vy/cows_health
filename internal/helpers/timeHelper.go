package helpers

import (
	"math"
)

// RoundTime функцию надо переписать
// она должна принимать duration  и возвразать
// возраст в виде троки: "XX лет XX месяцев"

func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}
