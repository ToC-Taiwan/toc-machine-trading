// Package utils package utils
package utils

import (
	"errors"

	"github.com/markcheno/go-talib"
)

// GenerareMAByCloseArr -.
func GenerareMAByCloseArr(closeArr []float64) (lastMa float64) {
	maArr := talib.Ma(closeArr, len(closeArr), talib.SMA)
	if len(maArr) == 0 {
		return 0
	}
	return maArr[len(maArr)-1]
}

// GetBiasRateByCloseArr -.
func GetBiasRateByCloseArr(closeArr []float64) (biasRate float64, err error) {
	ma := GenerareMAByCloseArr(closeArr)
	if ma == 0 {
		return 0, errors.New("no ma")
	}
	return Round(100*(closeArr[len(closeArr)-1]-ma)/ma, 2), err
}

// GenerateRSI -.
func GenerateRSI(input []float64) (rsi float64, err error) {
	rsiArr := talib.Rsi(input, len(input)-1)
	if len(rsiArr) == 0 {
		return 0, errors.New("no rsi")
	}
	return Round(rsiArr[len(rsiArr)-1], 2), err
}
