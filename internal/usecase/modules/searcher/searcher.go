// Package searcher package searcher
package searcher

import (
	"sort"
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

	stockArr  []*entity.Stock
	futureArr []*entity.Future
	optionArr []*entity.Option
}

func Get() Searcher {
	if singleton == nil {
		once.Do(func() {
			singleton = &searcher{}
		})
		return Get()
	}
	return singleton
}

func (s *searcher) AddStock(stock *entity.Stock) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.stockArr = append(s.stockArr, stock)
}

func (s *searcher) AddFuture(future *entity.Future) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if strings.Contains(strings.ToUpper(future.Code), "R1") || strings.Contains(strings.ToUpper(future.Code), "R2") {
		return
	}

	s.futureArr = append(s.futureArr, future)
}

func (s *searcher) AddOption(option *entity.Option) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.optionArr = append(s.optionArr, option)
}

func (s *searcher) SearchStock(param string) []*entity.Stock {
	if param == "" {
		return nil
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	param = strings.ToLower(param)
	var result []*entity.Stock
	for _, v := range s.stockArr {
		if strings.Contains(strings.ToLower(v.Number), param) {
			result = append(result, v)
		} else if strings.Contains(strings.ToLower(v.Name), param) {
			result = append(result, v)
		}
	}
	if len(result) != 0 {
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Number < result[j].Number
		})
	}
	return result
}

func (s *searcher) SearchFuture(param string) []*entity.Future {
	if param == "" {
		return nil
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	param = strings.ToLower(param)
	var result []*entity.Future
	for _, v := range s.futureArr {
		if strings.Contains(strings.ToLower(v.Code), param) {
			result = append(result, v)
		} else if strings.Contains(strings.ToLower(v.Name), param) {
			result = append(result, v)
		}
	}
	if len(result) != 0 {
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].DeliveryDate.Before(result[j].DeliveryDate)
		})
	}
	return result
}

func (s *searcher) SearchOption(param string) []*entity.Option {
	if param == "" {
		return nil
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	param = strings.ToLower(param)
	var result []*entity.Option
	for _, v := range s.optionArr {
		if strings.Contains(strings.ToLower(v.Code), param) {
			result = append(result, v)
		} else if strings.Contains(strings.ToLower(v.Name), param) {
			result = append(result, v)
		}
	}
	if len(result) != 0 {
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].DeliveryDate.Before(result[j].DeliveryDate)
		})
	}
	return result
}
