// Package utils package utils
package utils

import (
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
