// Package global package global
package global

import (
	"sync"
)

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

// Setting Setting
type Setting struct {
	lock          sync.RWMutex
	basePath      string
	isDevelopment bool
}

var globalSetting = &Setting{}

// SetBasePath SetBasePath
func SetBasePath(path string) {
	defer globalSetting.lock.RUnlock()
	globalSetting.lock.RLock()
	globalSetting.basePath = path
}

// GetBasePath GetBasePath
func GetBasePath() string {
	defer globalSetting.lock.RUnlock()
	globalSetting.lock.RLock()
	return globalSetting.basePath
}

// SetIsDevelopment SetIsDevelopment
func SetIsDevelopment(is bool) {
	defer globalSetting.lock.RUnlock()
	globalSetting.lock.RLock()
	globalSetting.isDevelopment = is
}

// GetIsDevelopment GetIsDevelopment
func GetIsDevelopment() bool {
	defer globalSetting.lock.RUnlock()
	globalSetting.lock.RLock()
	return globalSetting.isDevelopment
}
