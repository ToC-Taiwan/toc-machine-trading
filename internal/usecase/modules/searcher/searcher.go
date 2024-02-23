// Package searcher package searcher
package searcher

import (
	"strings"
	"sync"

	"tmt/internal/entity"
)

var (
	singleton Searcher
	once      sync.Once
)

type Searcher interface {
	AddStock(stock *entity.Stock)
	AddFuture(future *entity.Future)
	AddOption(option *entity.Option)

	SearchStock(code string) []*entity.Stock
	SearchFuture(code string) []*entity.Future
	SearchOption(code string) []*entity.Option
}

type searcher struct {
	lock sync.RWMutex

	stockMap     map[string]*entity.Stock
	stockCodeArr []string

	futureMap     map[string]*entity.Future
	futureCodeArr []string

	optionMap     map[string]*entity.Option
	optionCodeArr []string
}

func Get() Searcher {
	if singleton == nil {
		once.Do(func() {
			singleton = &searcher{
				stockMap:  make(map[string]*entity.Stock),
				futureMap: make(map[string]*entity.Future),
				optionMap: make(map[string]*entity.Option),
			}
		})
		return Get()
	}
	return singleton
}

func (s *searcher) AddStock(stock *entity.Stock) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.stockMap[stock.Code] = stock
	s.stockCodeArr = append(s.stockCodeArr, stock.Code)
}

func (s *searcher) AddFuture(future *entity.Future) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.futureMap[future.Code] = future
	s.futureCodeArr = append(s.futureCodeArr, future.Code)
}

func (s *searcher) AddOption(option *entity.Option) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.optionMap[option.Code] = option
	s.optionCodeArr = append(s.optionCodeArr, option.Code)
}

func (s *searcher) SearchStock(code string) []*entity.Stock {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var result []*entity.Stock
	for _, v := range s.stockCodeArr {
		if strings.Contains(v, code) {
			result = append(result, s.stockMap[v])
		}
	}
	return result
}

func (s *searcher) SearchFuture(code string) []*entity.Future {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var result []*entity.Future
	for _, v := range s.futureCodeArr {
		if strings.Contains(v, code) {
			result = append(result, s.futureMap[v])
		}
	}
	return result
}

func (s *searcher) SearchOption(code string) []*entity.Option {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var result []*entity.Option
	for _, v := range s.optionCodeArr {
		if strings.Contains(v, code) {
			result = append(result, s.optionMap[v])
		}
	}
	return result
}
