// Package utils package utils
package utils

import (
	"errors"

	"github.com/markcheno/go-talib"
)

// GenerareMAByCloseArr GenerareMAByCloseArr
func GenerareMAByCloseArr(closeArr []float64) (lastMa float64, err error) {
	maArr := talib.Ma(closeArr, len(closeArr), talib.SMA)
	if len(maArr) == 0 {
		return 0, errors.New("no ma")
	}
	return maArr[len(maArr)-1], err
}

// GetBiasRateByCloseArr GetBiasRateByCloseArr
func GetBiasRateByCloseArr(closeArr []float64) (biasRate float64, err error) {
	var ma float64
	ma, err = GenerareMAByCloseArr(closeArr)
	if err != nil {
		return biasRate, err
	}
	return Round(100*(closeArr[len(closeArr)-1]-ma)/ma, 2), err
}

// GenerateRSI GenerateRSI
func GenerateRSI(input []float64) (rsi float64, err error) {
	rsiArr := talib.Rsi(input, len(input)-1)
	if len(rsiArr) == 0 {
		return 0, errors.New("no rsi")
	}
	return Round(rsiArr[len(rsiArr)-1], 2), err
}
