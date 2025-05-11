package main

import (
	"fmt"
	"github.com/berserkkv/trader/service/connector"
)

// Function to calculate slope (m) and intercept (b) using the least squares method
func linearRegression(prices []float64) (float64, float64) {
	n := float64(len(prices))
	var sumX, sumY, sumXY, sumXX float64

	// Calculate the necessary sums
	for i, price := range prices {
		x := float64(i + 1) // Use 1-based index for x (starting from 1)
		y := price

		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	// Calculate the slope (m)
	m := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)

	// Calculate the intercept (b)
	b := (sumY - m*sumX) / n

	return m, b
}

// Function to calculate Linear Regression with offset
func calculateLinregWithOffset(prices []float64, length int, offset int) []float64 {
	var linreg []float64

	// Calculate the linreg values for the source series with the given length and offset
	for i := 0; i <= len(prices)-length; i++ {
		// Get the data points for the current window (length)
		window := prices[i : i+length]

		// Calculate the linear regression for the current window
		m, b := linearRegression(window)

		// Calculate the linreg value at the current point with the offset
		linregValue := b + m*float64(length-1-offset)

		// Store the linreg value
		linreg = append(linreg, linregValue)
	}

	return linreg
}

func main() {
	// Sample price data (e.g., close prices) in float64 format
	candles, _ := connector.FetchKlines("SOLUSDT", "15m", 200)

	prices := make([]float64, len(candles))
	for i, candle := range candles {
		prices[i] = candle.Low
	}

	// Define the length of the regression window
	length := 20

	// Define the offset (e.g., 8)
	offset := 7

	// Calculate linreg with offset
	linreg := calculateLinregWithOffset(prices, length, offset)

	// Output the linreg values
	fmt.Println("Linear Regression values with offset:")
	for i, val := range linreg {
		fmt.Printf("Period %d: %.2f\n", i+length, val)
	}
}
