// Package utils package utils
package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

// Round Round
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// RandomASCIILowerOctdigitsString -.
func RandomASCIILowerOctdigitsString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz01234567")

	s := make([]rune, n)
	for i := range s {
		randomBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		s[i] = letters[randomBigInt.Int64()]
	}
	return string(s)
}
