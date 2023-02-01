// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces_repo.go

// Package usecase_test is a generated GoMock package.
package usecase_test

import (
	context "context"
	reflect "reflect"
	time "time"
	entity "tmt/internal/entity"

	gomock "github.com/golang/mock/gomock"
)

// MockBasicRepo is a mock of BasicRepo interface.
type MockBasicRepo struct {
	ctrl     *gomock.Controller
	recorder *MockBasicRepoMockRecorder
}

// MockBasicRepoMockRecorder is the mock recorder for MockBasicRepo.
type MockBasicRepoMockRecorder struct {
	mock *MockBasicRepo
}

// NewMockBasicRepo creates a new mock instance.
func NewMockBasicRepo(ctrl *gomock.Controller) *MockBasicRepo {
	mock := &MockBasicRepo{ctrl: ctrl}
	mock.recorder = &MockBasicRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBasicRepo) EXPECT() *MockBasicRepoMockRecorder {
	return m.recorder
}

// InsertOrUpdatetCalendarDateArr mocks base method.
func (m *MockBasicRepo) InsertOrUpdatetCalendarDateArr(ctx context.Context, t []*entity.CalendarDate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdatetCalendarDateArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdatetCalendarDateArr indicates an expected call of InsertOrUpdatetCalendarDateArr.
func (mr *MockBasicRepoMockRecorder) InsertOrUpdatetCalendarDateArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdatetCalendarDateArr", reflect.TypeOf((*MockBasicRepo)(nil).InsertOrUpdatetCalendarDateArr), ctx, t)
}

// InsertOrUpdatetFutureArr mocks base method.
func (m *MockBasicRepo) InsertOrUpdatetFutureArr(ctx context.Context, t []*entity.Future) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdatetFutureArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdatetFutureArr indicates an expected call of InsertOrUpdatetFutureArr.
func (mr *MockBasicRepoMockRecorder) InsertOrUpdatetFutureArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdatetFutureArr", reflect.TypeOf((*MockBasicRepo)(nil).InsertOrUpdatetFutureArr), ctx, t)
}

// InsertOrUpdatetStockArr mocks base method.
func (m *MockBasicRepo) InsertOrUpdatetStockArr(ctx context.Context, t []*entity.Stock) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdatetStockArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdatetStockArr indicates an expected call of InsertOrUpdatetStockArr.
func (mr *MockBasicRepoMockRecorder) InsertOrUpdatetStockArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdatetStockArr", reflect.TypeOf((*MockBasicRepo)(nil).InsertOrUpdatetStockArr), ctx, t)
}

// QueryAllCalendar mocks base method.
func (m *MockBasicRepo) QueryAllCalendar(ctx context.Context) (map[time.Time]*entity.CalendarDate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllCalendar", ctx)
	ret0, _ := ret[0].(map[time.Time]*entity.CalendarDate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllCalendar indicates an expected call of QueryAllCalendar.
func (mr *MockBasicRepoMockRecorder) QueryAllCalendar(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllCalendar", reflect.TypeOf((*MockBasicRepo)(nil).QueryAllCalendar), ctx)
}

// QueryAllFuture mocks base method.
func (m *MockBasicRepo) QueryAllFuture(ctx context.Context) (map[string]*entity.Future, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllFuture", ctx)
	ret0, _ := ret[0].(map[string]*entity.Future)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllFuture indicates an expected call of QueryAllFuture.
func (mr *MockBasicRepoMockRecorder) QueryAllFuture(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllFuture", reflect.TypeOf((*MockBasicRepo)(nil).QueryAllFuture), ctx)
}

// QueryAllStock mocks base method.
func (m *MockBasicRepo) QueryAllStock(ctx context.Context) (map[string]*entity.Stock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllStock", ctx)
	ret0, _ := ret[0].(map[string]*entity.Stock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllStock indicates an expected call of QueryAllStock.
func (mr *MockBasicRepoMockRecorder) QueryAllStock(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllStock", reflect.TypeOf((*MockBasicRepo)(nil).QueryAllStock), ctx)
}

// UpdateAllStockDayTradeToNo mocks base method.
func (m *MockBasicRepo) UpdateAllStockDayTradeToNo(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAllStockDayTradeToNo", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAllStockDayTradeToNo indicates an expected call of UpdateAllStockDayTradeToNo.
func (mr *MockBasicRepoMockRecorder) UpdateAllStockDayTradeToNo(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAllStockDayTradeToNo", reflect.TypeOf((*MockBasicRepo)(nil).UpdateAllStockDayTradeToNo), ctx)
}

// MockTargetRepo is a mock of TargetRepo interface.
type MockTargetRepo struct {
	ctrl     *gomock.Controller
	recorder *MockTargetRepoMockRecorder
}

// MockTargetRepoMockRecorder is the mock recorder for MockTargetRepo.
type MockTargetRepoMockRecorder struct {
	mock *MockTargetRepo
}

// NewMockTargetRepo creates a new mock instance.
func NewMockTargetRepo(ctrl *gomock.Controller) *MockTargetRepo {
	mock := &MockTargetRepo{ctrl: ctrl}
	mock.recorder = &MockTargetRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTargetRepo) EXPECT() *MockTargetRepoMockRecorder {
	return m.recorder
}

// InsertOrUpdateTargetArr mocks base method.
func (m *MockTargetRepo) InsertOrUpdateTargetArr(ctx context.Context, t []*entity.StockTarget) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdateTargetArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdateTargetArr indicates an expected call of InsertOrUpdateTargetArr.
func (mr *MockTargetRepoMockRecorder) InsertOrUpdateTargetArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdateTargetArr", reflect.TypeOf((*MockTargetRepo)(nil).InsertOrUpdateTargetArr), ctx, t)
}

// QueryAllMXFFuture mocks base method.
func (m *MockTargetRepo) QueryAllMXFFuture(ctx context.Context) ([]*entity.Future, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllMXFFuture", ctx)
	ret0, _ := ret[0].([]*entity.Future)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllMXFFuture indicates an expected call of QueryAllMXFFuture.
func (mr *MockTargetRepoMockRecorder) QueryAllMXFFuture(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllMXFFuture", reflect.TypeOf((*MockTargetRepo)(nil).QueryAllMXFFuture), ctx)
}

// QueryTargetsByTradeDay mocks base method.
func (m *MockTargetRepo) QueryTargetsByTradeDay(ctx context.Context, tradeDay time.Time) ([]*entity.StockTarget, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTargetsByTradeDay", ctx, tradeDay)
	ret0, _ := ret[0].([]*entity.StockTarget)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryTargetsByTradeDay indicates an expected call of QueryTargetsByTradeDay.
func (mr *MockTargetRepoMockRecorder) QueryTargetsByTradeDay(ctx, tradeDay interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTargetsByTradeDay", reflect.TypeOf((*MockTargetRepo)(nil).QueryTargetsByTradeDay), ctx, tradeDay)
}

// MockHistoryRepo is a mock of HistoryRepo interface.
type MockHistoryRepo struct {
	ctrl     *gomock.Controller
	recorder *MockHistoryRepoMockRecorder
}

// MockHistoryRepoMockRecorder is the mock recorder for MockHistoryRepo.
type MockHistoryRepoMockRecorder struct {
	mock *MockHistoryRepo
}

// NewMockHistoryRepo creates a new mock instance.
func NewMockHistoryRepo(ctrl *gomock.Controller) *MockHistoryRepo {
	mock := &MockHistoryRepo{ctrl: ctrl}
	mock.recorder = &MockHistoryRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistoryRepo) EXPECT() *MockHistoryRepoMockRecorder {
	return m.recorder
}

// DeleteHistoryCloseByStockAndDate mocks base method.
func (m *MockHistoryRepo) DeleteHistoryCloseByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteHistoryCloseByStockAndDate", ctx, stockNumArr, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHistoryCloseByStockAndDate indicates an expected call of DeleteHistoryCloseByStockAndDate.
func (mr *MockHistoryRepoMockRecorder) DeleteHistoryCloseByStockAndDate(ctx, stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHistoryCloseByStockAndDate", reflect.TypeOf((*MockHistoryRepo)(nil).DeleteHistoryCloseByStockAndDate), ctx, stockNumArr, date)
}

// DeleteHistoryKbarByStockAndDate mocks base method.
func (m *MockHistoryRepo) DeleteHistoryKbarByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteHistoryKbarByStockAndDate", ctx, stockNumArr, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHistoryKbarByStockAndDate indicates an expected call of DeleteHistoryKbarByStockAndDate.
func (mr *MockHistoryRepoMockRecorder) DeleteHistoryKbarByStockAndDate(ctx, stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHistoryKbarByStockAndDate", reflect.TypeOf((*MockHistoryRepo)(nil).DeleteHistoryKbarByStockAndDate), ctx, stockNumArr, date)
}

// DeleteHistoryTickByStockAndDate mocks base method.
func (m *MockHistoryRepo) DeleteHistoryTickByStockAndDate(ctx context.Context, stockNumArr []string, date time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteHistoryTickByStockAndDate", ctx, stockNumArr, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHistoryTickByStockAndDate indicates an expected call of DeleteHistoryTickByStockAndDate.
func (mr *MockHistoryRepoMockRecorder) DeleteHistoryTickByStockAndDate(ctx, stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHistoryTickByStockAndDate", reflect.TypeOf((*MockHistoryRepo)(nil).DeleteHistoryTickByStockAndDate), ctx, stockNumArr, date)
}

// InsertFutureHistoryClose mocks base method.
func (m *MockHistoryRepo) InsertFutureHistoryClose(ctx context.Context, c *entity.FutureHistoryClose) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertFutureHistoryClose", ctx, c)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertFutureHistoryClose indicates an expected call of InsertFutureHistoryClose.
func (mr *MockHistoryRepoMockRecorder) InsertFutureHistoryClose(ctx, c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertFutureHistoryClose", reflect.TypeOf((*MockHistoryRepo)(nil).InsertFutureHistoryClose), ctx, c)
}

// InsertFutureHistoryTickArr mocks base method.
func (m *MockHistoryRepo) InsertFutureHistoryTickArr(ctx context.Context, t []*entity.FutureHistoryTick) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertFutureHistoryTickArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertFutureHistoryTickArr indicates an expected call of InsertFutureHistoryTickArr.
func (mr *MockHistoryRepoMockRecorder) InsertFutureHistoryTickArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertFutureHistoryTickArr", reflect.TypeOf((*MockHistoryRepo)(nil).InsertFutureHistoryTickArr), ctx, t)
}

// InsertHistoryCloseArr mocks base method.
func (m *MockHistoryRepo) InsertHistoryCloseArr(ctx context.Context, t []*entity.StockHistoryClose) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertHistoryCloseArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertHistoryCloseArr indicates an expected call of InsertHistoryCloseArr.
func (mr *MockHistoryRepoMockRecorder) InsertHistoryCloseArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertHistoryCloseArr", reflect.TypeOf((*MockHistoryRepo)(nil).InsertHistoryCloseArr), ctx, t)
}

// InsertHistoryKbarArr mocks base method.
func (m *MockHistoryRepo) InsertHistoryKbarArr(ctx context.Context, t []*entity.StockHistoryKbar) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertHistoryKbarArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertHistoryKbarArr indicates an expected call of InsertHistoryKbarArr.
func (mr *MockHistoryRepoMockRecorder) InsertHistoryKbarArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertHistoryKbarArr", reflect.TypeOf((*MockHistoryRepo)(nil).InsertHistoryKbarArr), ctx, t)
}

// InsertHistoryTickArr mocks base method.
func (m *MockHistoryRepo) InsertHistoryTickArr(ctx context.Context, t []*entity.StockHistoryTick) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertHistoryTickArr", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertHistoryTickArr indicates an expected call of InsertHistoryTickArr.
func (mr *MockHistoryRepoMockRecorder) InsertHistoryTickArr(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertHistoryTickArr", reflect.TypeOf((*MockHistoryRepo)(nil).InsertHistoryTickArr), ctx, t)
}

// InsertQuaterMA mocks base method.
func (m *MockHistoryRepo) InsertQuaterMA(ctx context.Context, t *entity.StockHistoryAnalyze) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertQuaterMA", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertQuaterMA indicates an expected call of InsertQuaterMA.
func (mr *MockHistoryRepoMockRecorder) InsertQuaterMA(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertQuaterMA", reflect.TypeOf((*MockHistoryRepo)(nil).InsertQuaterMA), ctx, t)
}

// QueryAllQuaterMAByStockNum mocks base method.
func (m *MockHistoryRepo) QueryAllQuaterMAByStockNum(ctx context.Context, stockNum string) (map[time.Time]*entity.StockHistoryAnalyze, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllQuaterMAByStockNum", ctx, stockNum)
	ret0, _ := ret[0].(map[time.Time]*entity.StockHistoryAnalyze)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllQuaterMAByStockNum indicates an expected call of QueryAllQuaterMAByStockNum.
func (mr *MockHistoryRepoMockRecorder) QueryAllQuaterMAByStockNum(ctx, stockNum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllQuaterMAByStockNum", reflect.TypeOf((*MockHistoryRepo)(nil).QueryAllQuaterMAByStockNum), ctx, stockNum)
}

// QueryFutureHistoryCloseByDate mocks base method.
func (m *MockHistoryRepo) QueryFutureHistoryCloseByDate(ctx context.Context, code string, tradeDay time.Time) (*entity.FutureHistoryClose, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryFutureHistoryCloseByDate", ctx, code, tradeDay)
	ret0, _ := ret[0].(*entity.FutureHistoryClose)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryFutureHistoryCloseByDate indicates an expected call of QueryFutureHistoryCloseByDate.
func (mr *MockHistoryRepoMockRecorder) QueryFutureHistoryCloseByDate(ctx, code, tradeDay interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryFutureHistoryCloseByDate", reflect.TypeOf((*MockHistoryRepo)(nil).QueryFutureHistoryCloseByDate), ctx, code, tradeDay)
}

// QueryFutureHistoryTickArrByTime mocks base method.
func (m *MockHistoryRepo) QueryFutureHistoryTickArrByTime(ctx context.Context, code string, startTime, endTime time.Time) ([]*entity.FutureHistoryTick, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryFutureHistoryTickArrByTime", ctx, code, startTime, endTime)
	ret0, _ := ret[0].([]*entity.FutureHistoryTick)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryFutureHistoryTickArrByTime indicates an expected call of QueryFutureHistoryTickArrByTime.
func (mr *MockHistoryRepoMockRecorder) QueryFutureHistoryTickArrByTime(ctx, code, startTime, endTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryFutureHistoryTickArrByTime", reflect.TypeOf((*MockHistoryRepo)(nil).QueryFutureHistoryTickArrByTime), ctx, code, startTime, endTime)
}

// QueryMultiStockKbarArrByDate mocks base method.
func (m *MockHistoryRepo) QueryMultiStockKbarArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryKbar, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryMultiStockKbarArrByDate", ctx, stockNumArr, date)
	ret0, _ := ret[0].(map[string][]*entity.StockHistoryKbar)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryMultiStockKbarArrByDate indicates an expected call of QueryMultiStockKbarArrByDate.
func (mr *MockHistoryRepoMockRecorder) QueryMultiStockKbarArrByDate(ctx, stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryMultiStockKbarArrByDate", reflect.TypeOf((*MockHistoryRepo)(nil).QueryMultiStockKbarArrByDate), ctx, stockNumArr, date)
}

// QueryMultiStockTickArrByDate mocks base method.
func (m *MockHistoryRepo) QueryMultiStockTickArrByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string][]*entity.StockHistoryTick, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryMultiStockTickArrByDate", ctx, stockNumArr, date)
	ret0, _ := ret[0].(map[string][]*entity.StockHistoryTick)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryMultiStockTickArrByDate indicates an expected call of QueryMultiStockTickArrByDate.
func (mr *MockHistoryRepoMockRecorder) QueryMultiStockTickArrByDate(ctx, stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryMultiStockTickArrByDate", reflect.TypeOf((*MockHistoryRepo)(nil).QueryMultiStockTickArrByDate), ctx, stockNumArr, date)
}

// QueryMutltiStockCloseByDate mocks base method.
func (m *MockHistoryRepo) QueryMutltiStockCloseByDate(ctx context.Context, stockNumArr []string, date time.Time) (map[string]*entity.StockHistoryClose, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryMutltiStockCloseByDate", ctx, stockNumArr, date)
	ret0, _ := ret[0].(map[string]*entity.StockHistoryClose)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryMutltiStockCloseByDate indicates an expected call of QueryMutltiStockCloseByDate.
func (mr *MockHistoryRepoMockRecorder) QueryMutltiStockCloseByDate(ctx, stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryMutltiStockCloseByDate", reflect.TypeOf((*MockHistoryRepo)(nil).QueryMutltiStockCloseByDate), ctx, stockNumArr, date)
}

// MockRealTimeRepo is a mock of RealTimeRepo interface.
type MockRealTimeRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRealTimeRepoMockRecorder
}

// MockRealTimeRepoMockRecorder is the mock recorder for MockRealTimeRepo.
type MockRealTimeRepoMockRecorder struct {
	mock *MockRealTimeRepo
}

// NewMockRealTimeRepo creates a new mock instance.
func NewMockRealTimeRepo(ctrl *gomock.Controller) *MockRealTimeRepo {
	mock := &MockRealTimeRepo{ctrl: ctrl}
	mock.recorder = &MockRealTimeRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRealTimeRepo) EXPECT() *MockRealTimeRepoMockRecorder {
	return m.recorder
}

// InsertEvent mocks base method.
func (m *MockRealTimeRepo) InsertEvent(ctx context.Context, t *entity.SinopacEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertEvent", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertEvent indicates an expected call of InsertEvent.
func (mr *MockRealTimeRepoMockRecorder) InsertEvent(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertEvent", reflect.TypeOf((*MockRealTimeRepo)(nil).InsertEvent), ctx, t)
}

// MockTradeRepo is a mock of TradeRepo interface.
type MockTradeRepo struct {
	ctrl     *gomock.Controller
	recorder *MockTradeRepoMockRecorder
}

// MockTradeRepoMockRecorder is the mock recorder for MockTradeRepo.
type MockTradeRepoMockRecorder struct {
	mock *MockTradeRepo
}

// NewMockTradeRepo creates a new mock instance.
func NewMockTradeRepo(ctrl *gomock.Controller) *MockTradeRepo {
	mock := &MockTradeRepo{ctrl: ctrl}
	mock.recorder = &MockTradeRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTradeRepo) EXPECT() *MockTradeRepoMockRecorder {
	return m.recorder
}

// InsertOrUpdateFutureOrderByOrderID mocks base method.
func (m *MockTradeRepo) InsertOrUpdateFutureOrderByOrderID(ctx context.Context, t *entity.FutureOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdateFutureOrderByOrderID", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdateFutureOrderByOrderID indicates an expected call of InsertOrUpdateFutureOrderByOrderID.
func (mr *MockTradeRepoMockRecorder) InsertOrUpdateFutureOrderByOrderID(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdateFutureOrderByOrderID", reflect.TypeOf((*MockTradeRepo)(nil).InsertOrUpdateFutureOrderByOrderID), ctx, t)
}

// InsertOrUpdateFutureTradeBalance mocks base method.
func (m *MockTradeRepo) InsertOrUpdateFutureTradeBalance(ctx context.Context, t *entity.FutureTradeBalance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdateFutureTradeBalance", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdateFutureTradeBalance indicates an expected call of InsertOrUpdateFutureTradeBalance.
func (mr *MockTradeRepoMockRecorder) InsertOrUpdateFutureTradeBalance(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdateFutureTradeBalance", reflect.TypeOf((*MockTradeRepo)(nil).InsertOrUpdateFutureTradeBalance), ctx, t)
}

// InsertOrUpdateOrderByOrderID mocks base method.
func (m *MockTradeRepo) InsertOrUpdateOrderByOrderID(ctx context.Context, t *entity.StockOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdateOrderByOrderID", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdateOrderByOrderID indicates an expected call of InsertOrUpdateOrderByOrderID.
func (mr *MockTradeRepoMockRecorder) InsertOrUpdateOrderByOrderID(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdateOrderByOrderID", reflect.TypeOf((*MockTradeRepo)(nil).InsertOrUpdateOrderByOrderID), ctx, t)
}

// InsertOrUpdateStockTradeBalance mocks base method.
func (m *MockTradeRepo) InsertOrUpdateStockTradeBalance(ctx context.Context, t *entity.StockTradeBalance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrUpdateStockTradeBalance", ctx, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrUpdateStockTradeBalance indicates an expected call of InsertOrUpdateStockTradeBalance.
func (mr *MockTradeRepoMockRecorder) InsertOrUpdateStockTradeBalance(ctx, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrUpdateStockTradeBalance", reflect.TypeOf((*MockTradeRepo)(nil).InsertOrUpdateStockTradeBalance), ctx, t)
}

// QueryAllFutureOrder mocks base method.
func (m *MockTradeRepo) QueryAllFutureOrder(ctx context.Context) ([]*entity.FutureOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllFutureOrder", ctx)
	ret0, _ := ret[0].([]*entity.FutureOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllFutureOrder indicates an expected call of QueryAllFutureOrder.
func (mr *MockTradeRepoMockRecorder) QueryAllFutureOrder(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllFutureOrder", reflect.TypeOf((*MockTradeRepo)(nil).QueryAllFutureOrder), ctx)
}

// QueryAllFutureOrderByDate mocks base method.
func (m *MockTradeRepo) QueryAllFutureOrderByDate(ctx context.Context, timeTange []time.Time) ([]*entity.FutureOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllFutureOrderByDate", ctx, timeTange)
	ret0, _ := ret[0].([]*entity.FutureOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllFutureOrderByDate indicates an expected call of QueryAllFutureOrderByDate.
func (mr *MockTradeRepoMockRecorder) QueryAllFutureOrderByDate(ctx, timeTange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllFutureOrderByDate", reflect.TypeOf((*MockTradeRepo)(nil).QueryAllFutureOrderByDate), ctx, timeTange)
}

// QueryAllFutureTradeBalance mocks base method.
func (m *MockTradeRepo) QueryAllFutureTradeBalance(ctx context.Context) ([]*entity.FutureTradeBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllFutureTradeBalance", ctx)
	ret0, _ := ret[0].([]*entity.FutureTradeBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllFutureTradeBalance indicates an expected call of QueryAllFutureTradeBalance.
func (mr *MockTradeRepoMockRecorder) QueryAllFutureTradeBalance(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllFutureTradeBalance", reflect.TypeOf((*MockTradeRepo)(nil).QueryAllFutureTradeBalance), ctx)
}

// QueryAllStockOrder mocks base method.
func (m *MockTradeRepo) QueryAllStockOrder(ctx context.Context) ([]*entity.StockOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllStockOrder", ctx)
	ret0, _ := ret[0].([]*entity.StockOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllStockOrder indicates an expected call of QueryAllStockOrder.
func (mr *MockTradeRepoMockRecorder) QueryAllStockOrder(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllStockOrder", reflect.TypeOf((*MockTradeRepo)(nil).QueryAllStockOrder), ctx)
}

// QueryAllStockOrderByDate mocks base method.
func (m *MockTradeRepo) QueryAllStockOrderByDate(ctx context.Context, timeTange []time.Time) ([]*entity.StockOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllStockOrderByDate", ctx, timeTange)
	ret0, _ := ret[0].([]*entity.StockOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllStockOrderByDate indicates an expected call of QueryAllStockOrderByDate.
func (mr *MockTradeRepoMockRecorder) QueryAllStockOrderByDate(ctx, timeTange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllStockOrderByDate", reflect.TypeOf((*MockTradeRepo)(nil).QueryAllStockOrderByDate), ctx, timeTange)
}

// QueryAllStockTradeBalance mocks base method.
func (m *MockTradeRepo) QueryAllStockTradeBalance(ctx context.Context) ([]*entity.StockTradeBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllStockTradeBalance", ctx)
	ret0, _ := ret[0].([]*entity.StockTradeBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllStockTradeBalance indicates an expected call of QueryAllStockTradeBalance.
func (mr *MockTradeRepoMockRecorder) QueryAllStockTradeBalance(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllStockTradeBalance", reflect.TypeOf((*MockTradeRepo)(nil).QueryAllStockTradeBalance), ctx)
}

// QueryFutureOrderByID mocks base method.
func (m *MockTradeRepo) QueryFutureOrderByID(ctx context.Context, orderID string) (*entity.FutureOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryFutureOrderByID", ctx, orderID)
	ret0, _ := ret[0].(*entity.FutureOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryFutureOrderByID indicates an expected call of QueryFutureOrderByID.
func (mr *MockTradeRepoMockRecorder) QueryFutureOrderByID(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryFutureOrderByID", reflect.TypeOf((*MockTradeRepo)(nil).QueryFutureOrderByID), ctx, orderID)
}

// QueryFutureTradeBalanceByDate mocks base method.
func (m *MockTradeRepo) QueryFutureTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.FutureTradeBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryFutureTradeBalanceByDate", ctx, date)
	ret0, _ := ret[0].(*entity.FutureTradeBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryFutureTradeBalanceByDate indicates an expected call of QueryFutureTradeBalanceByDate.
func (mr *MockTradeRepoMockRecorder) QueryFutureTradeBalanceByDate(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryFutureTradeBalanceByDate", reflect.TypeOf((*MockTradeRepo)(nil).QueryFutureTradeBalanceByDate), ctx, date)
}

// QueryLastFutureTradeBalance mocks base method.
func (m *MockTradeRepo) QueryLastFutureTradeBalance(ctx context.Context) (*entity.FutureTradeBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLastFutureTradeBalance", ctx)
	ret0, _ := ret[0].(*entity.FutureTradeBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryLastFutureTradeBalance indicates an expected call of QueryLastFutureTradeBalance.
func (mr *MockTradeRepoMockRecorder) QueryLastFutureTradeBalance(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLastFutureTradeBalance", reflect.TypeOf((*MockTradeRepo)(nil).QueryLastFutureTradeBalance), ctx)
}

// QueryLastStockTradeBalance mocks base method.
func (m *MockTradeRepo) QueryLastStockTradeBalance(ctx context.Context) (*entity.StockTradeBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLastStockTradeBalance", ctx)
	ret0, _ := ret[0].(*entity.StockTradeBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryLastStockTradeBalance indicates an expected call of QueryLastStockTradeBalance.
func (mr *MockTradeRepoMockRecorder) QueryLastStockTradeBalance(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLastStockTradeBalance", reflect.TypeOf((*MockTradeRepo)(nil).QueryLastStockTradeBalance), ctx)
}

// QueryStockOrderByID mocks base method.
func (m *MockTradeRepo) QueryStockOrderByID(ctx context.Context, orderID string) (*entity.StockOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryStockOrderByID", ctx, orderID)
	ret0, _ := ret[0].(*entity.StockOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryStockOrderByID indicates an expected call of QueryStockOrderByID.
func (mr *MockTradeRepoMockRecorder) QueryStockOrderByID(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryStockOrderByID", reflect.TypeOf((*MockTradeRepo)(nil).QueryStockOrderByID), ctx, orderID)
}

// QueryStockTradeBalanceByDate mocks base method.
func (m *MockTradeRepo) QueryStockTradeBalanceByDate(ctx context.Context, date time.Time) (*entity.StockTradeBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryStockTradeBalanceByDate", ctx, date)
	ret0, _ := ret[0].(*entity.StockTradeBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryStockTradeBalanceByDate indicates an expected call of QueryStockTradeBalanceByDate.
func (mr *MockTradeRepoMockRecorder) QueryStockTradeBalanceByDate(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryStockTradeBalanceByDate", reflect.TypeOf((*MockTradeRepo)(nil).QueryStockTradeBalanceByDate), ctx, date)
}