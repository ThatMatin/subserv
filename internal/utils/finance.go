package utils

import "math"

func CalculateFinalAmount(priceCents int, taxRate uint8) int {
	total := float64(priceCents) * (1 + float64(taxRate/100.))
	return int(math.Round(total))
}
