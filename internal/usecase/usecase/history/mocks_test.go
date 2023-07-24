// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package history_test is a generated GoMock package.
package history_test

import (
	context "context"
	reflect "reflect"
	time "time"
	config "tmt/cmd/config"
	entity "tmt/internal/entity"
	cache "tmt/internal/usecase/cache"
	simulator "tmt/internal/usecase/module/simulator"
	history "tmt/internal/usecase/usecase/history"
	pb "tmt/pb"
	eventbus "tmt/pkg/eventbus"
	log "tmt/pkg/log"

	gomock "github.com/golang/mock/gomock"
)

// MockHistory is a mock of History interface.
type MockHistory struct {
	ctrl     *gomock.Controller
	recorder *MockHistoryMockRecorder
}

// MockHistoryMockRecorder is the mock recorder for MockHistory.
type MockHistoryMockRecorder struct {
	mock *MockHistory
}

// NewMockHistory creates a new mock instance.
func NewMockHistory(ctrl *gomock.Controller) *MockHistory {
	mock := &MockHistory{ctrl: ctrl}
	mock.recorder = &MockHistoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistory) EXPECT() *MockHistoryMockRecorder {
	return m.recorder
}

// FetchFutureHistoryKbar mocks base method.
func (m *MockHistory) FetchFutureHistoryKbar(code string, date time.Time) ([]*entity.FutureHistoryKbar, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchFutureHistoryKbar", code, date)
	ret0, _ := ret[0].([]*entity.FutureHistoryKbar)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchFutureHistoryKbar indicates an expected call of FetchFutureHistoryKbar.
func (mr *MockHistoryMockRecorder) FetchFutureHistoryKbar(code, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchFutureHistoryKbar", reflect.TypeOf((*MockHistory)(nil).FetchFutureHistoryKbar), code, date)
}

// GetDayKbarByStockNumDate mocks base method.
func (m *MockHistory) GetDayKbarByStockNumDate(stockNum string, date time.Time) *entity.StockHistoryKbar {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDayKbarByStockNumDate", stockNum, date)
	ret0, _ := ret[0].(*entity.StockHistoryKbar)
	return ret0
}

// GetDayKbarByStockNumDate indicates an expected call of GetDayKbarByStockNumDate.
func (mr *MockHistoryMockRecorder) GetDayKbarByStockNumDate(stockNum, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDayKbarByStockNumDate", reflect.TypeOf((*MockHistory)(nil).GetDayKbarByStockNumDate), stockNum, date)
}

// GetTradeDay mocks base method.
func (m *MockHistory) GetTradeDay() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTradeDay")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetTradeDay indicates an expected call of GetTradeDay.
func (mr *MockHistoryMockRecorder) GetTradeDay() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTradeDay", reflect.TypeOf((*MockHistory)(nil).GetTradeDay))
}

// Init mocks base method.
func (m *MockHistory) Init(logger *log.Log, cc *cache.Cache, bus *eventbus.Bus) history.History {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", logger, cc, bus)
	ret0, _ := ret[0].(history.History)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockHistoryMockRecorder) Init(logger, cc, bus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockHistory)(nil).Init), logger, cc, bus)
}

// SimulateMulti mocks base method.
func (m *MockHistory) SimulateMulti(cond []*config.TradeFuture) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SimulateMulti", cond)
}

// SimulateMulti indicates an expected call of SimulateMulti.
func (mr *MockHistoryMockRecorder) SimulateMulti(cond interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SimulateMulti", reflect.TypeOf((*MockHistory)(nil).SimulateMulti), cond)
}

// SimulateOne mocks base method.
func (m *MockHistory) SimulateOne(cond *config.TradeFuture) *simulator.SimulateBalance {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SimulateOne", cond)
	ret0, _ := ret[0].(*simulator.SimulateBalance)
	return ret0
}

// SimulateOne indicates an expected call of SimulateOne.
func (mr *MockHistoryMockRecorder) SimulateOne(cond interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SimulateOne", reflect.TypeOf((*MockHistory)(nil).SimulateOne), cond)
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

// MockHistorygRPCAPI is a mock of HistorygRPCAPI interface.
type MockHistorygRPCAPI struct {
	ctrl     *gomock.Controller
	recorder *MockHistorygRPCAPIMockRecorder
}

// MockHistorygRPCAPIMockRecorder is the mock recorder for MockHistorygRPCAPI.
type MockHistorygRPCAPIMockRecorder struct {
	mock *MockHistorygRPCAPI
}

// NewMockHistorygRPCAPI creates a new mock instance.
func NewMockHistorygRPCAPI(ctrl *gomock.Controller) *MockHistorygRPCAPI {
	mock := &MockHistorygRPCAPI{ctrl: ctrl}
	mock.recorder = &MockHistorygRPCAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistorygRPCAPI) EXPECT() *MockHistorygRPCAPIMockRecorder {
	return m.recorder
}

// GetFutureHistoryClose mocks base method.
func (m *MockHistorygRPCAPI) GetFutureHistoryClose(codeArr []string, date string) ([]*pb.HistoryCloseMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFutureHistoryClose", codeArr, date)
	ret0, _ := ret[0].([]*pb.HistoryCloseMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFutureHistoryClose indicates an expected call of GetFutureHistoryClose.
func (mr *MockHistorygRPCAPIMockRecorder) GetFutureHistoryClose(codeArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFutureHistoryClose", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetFutureHistoryClose), codeArr, date)
}

// GetFutureHistoryKbar mocks base method.
func (m *MockHistorygRPCAPI) GetFutureHistoryKbar(codeArr []string, date string) ([]*pb.HistoryKbarMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFutureHistoryKbar", codeArr, date)
	ret0, _ := ret[0].([]*pb.HistoryKbarMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFutureHistoryKbar indicates an expected call of GetFutureHistoryKbar.
func (mr *MockHistorygRPCAPIMockRecorder) GetFutureHistoryKbar(codeArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFutureHistoryKbar", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetFutureHistoryKbar), codeArr, date)
}

// GetFutureHistoryTick mocks base method.
func (m *MockHistorygRPCAPI) GetFutureHistoryTick(codeArr []string, date string) ([]*pb.HistoryTickMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFutureHistoryTick", codeArr, date)
	ret0, _ := ret[0].([]*pb.HistoryTickMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFutureHistoryTick indicates an expected call of GetFutureHistoryTick.
func (mr *MockHistorygRPCAPIMockRecorder) GetFutureHistoryTick(codeArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFutureHistoryTick", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetFutureHistoryTick), codeArr, date)
}

// GetStockHistoryClose mocks base method.
func (m *MockHistorygRPCAPI) GetStockHistoryClose(stockNumArr []string, date string) ([]*pb.HistoryCloseMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStockHistoryClose", stockNumArr, date)
	ret0, _ := ret[0].([]*pb.HistoryCloseMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStockHistoryClose indicates an expected call of GetStockHistoryClose.
func (mr *MockHistorygRPCAPIMockRecorder) GetStockHistoryClose(stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStockHistoryClose", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetStockHistoryClose), stockNumArr, date)
}

// GetStockHistoryCloseByDateArr mocks base method.
func (m *MockHistorygRPCAPI) GetStockHistoryCloseByDateArr(stockNumArr, date []string) ([]*pb.HistoryCloseMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStockHistoryCloseByDateArr", stockNumArr, date)
	ret0, _ := ret[0].([]*pb.HistoryCloseMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStockHistoryCloseByDateArr indicates an expected call of GetStockHistoryCloseByDateArr.
func (mr *MockHistorygRPCAPIMockRecorder) GetStockHistoryCloseByDateArr(stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStockHistoryCloseByDateArr", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetStockHistoryCloseByDateArr), stockNumArr, date)
}

// GetStockHistoryKbar mocks base method.
func (m *MockHistorygRPCAPI) GetStockHistoryKbar(stockNumArr []string, date string) ([]*pb.HistoryKbarMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStockHistoryKbar", stockNumArr, date)
	ret0, _ := ret[0].([]*pb.HistoryKbarMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStockHistoryKbar indicates an expected call of GetStockHistoryKbar.
func (mr *MockHistorygRPCAPIMockRecorder) GetStockHistoryKbar(stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStockHistoryKbar", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetStockHistoryKbar), stockNumArr, date)
}

// GetStockHistoryTick mocks base method.
func (m *MockHistorygRPCAPI) GetStockHistoryTick(stockNumArr []string, date string) ([]*pb.HistoryTickMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStockHistoryTick", stockNumArr, date)
	ret0, _ := ret[0].([]*pb.HistoryTickMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStockHistoryTick indicates an expected call of GetStockHistoryTick.
func (mr *MockHistorygRPCAPIMockRecorder) GetStockHistoryTick(stockNumArr, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStockHistoryTick", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetStockHistoryTick), stockNumArr, date)
}

// GetStockTSEHistoryClose mocks base method.
func (m *MockHistorygRPCAPI) GetStockTSEHistoryClose(date string) ([]*pb.HistoryCloseMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStockTSEHistoryClose", date)
	ret0, _ := ret[0].([]*pb.HistoryCloseMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStockTSEHistoryClose indicates an expected call of GetStockTSEHistoryClose.
func (mr *MockHistorygRPCAPIMockRecorder) GetStockTSEHistoryClose(date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStockTSEHistoryClose", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetStockTSEHistoryClose), date)
}

// GetStockTSEHistoryKbar mocks base method.
func (m *MockHistorygRPCAPI) GetStockTSEHistoryKbar(date string) ([]*pb.HistoryKbarMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStockTSEHistoryKbar", date)
	ret0, _ := ret[0].([]*pb.HistoryKbarMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStockTSEHistoryKbar indicates an expected call of GetStockTSEHistoryKbar.
func (mr *MockHistorygRPCAPIMockRecorder) GetStockTSEHistoryKbar(date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStockTSEHistoryKbar", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetStockTSEHistoryKbar), date)
}

// GetStockTSEHistoryTick mocks base method.
func (m *MockHistorygRPCAPI) GetStockTSEHistoryTick(date string) ([]*pb.HistoryTickMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStockTSEHistoryTick", date)
	ret0, _ := ret[0].([]*pb.HistoryTickMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStockTSEHistoryTick indicates an expected call of GetStockTSEHistoryTick.
func (mr *MockHistorygRPCAPIMockRecorder) GetStockTSEHistoryTick(date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStockTSEHistoryTick", reflect.TypeOf((*MockHistorygRPCAPI)(nil).GetStockTSEHistoryTick), date)
}