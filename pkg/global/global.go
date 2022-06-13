// Package global package global
package global

const (
	// LongTimeLayout LongTimeLayout
	LongTimeLayout string = "2006-01-02 15:04:05"
	// ShortTimeLayout ShortTimeLayout
	ShortTimeLayout string = "2006-01-02"
)

const (
	// StartTradeYear -.
	StartTradeYear int = 2021
	// EndTradeYear -.
	EndTradeYear int = 2022
)

var basePath string

// SetBasePath SetBasePath
func SetBasePath(path string) {
	basePath = path
}

// GetBasePath GetBasePath
func GetBasePath() string {
	return basePath
}
