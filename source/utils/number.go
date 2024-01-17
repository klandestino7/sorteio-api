package utils

import (
	"fmt"
	"math"
)

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(math.Round(num*output)) / output
}

func PriceToFloatString(price uint32) string {
	var fPrice float32 = float32(price) / 100.0
	// fmt.Printf("%T = %d, %T = %f\n", price, price, fPrice, fPrice)
	stringPrice := fmt.Sprintf("%.2f", fPrice)

	DebugPrint("PriceToFloatString :: ", stringPrice)

	return stringPrice
}
