// Package utils package utils
package utils

import (
	"errors"
	"math"

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
func GenerateRSI(input []float64, effTimes int) float64 {
	baseClose := input[0]
	diff := GetStockDiffByClose(baseClose)

	var positive, negative float64
	for _, v := range input[1:] {
		gap := v - baseClose
		time := math.Abs(Round(gap/diff, 0))
		switch {
		case gap > 0:
			positive += time
		case gap < 0:
			negative += time
		}

		if gap != 0 {
			baseClose = v
			diff = GetStockDiffByClose(baseClose)
		}
	}

	if totalEff := positive + negative; totalEff < float64(effTimes) {
		return 0
	}

	return Round(100*positive/(positive+negative), 2)
}
