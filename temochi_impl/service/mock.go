// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	repository "github.com/bagusbpg/tenpo/temochi_impl/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// DeleteChannelStock mocks base method.
func (m *MockRepository) DeleteChannelStock(ctx context.Context, input repository.DeleteChannelStockDBInput, output *repository.DeleteChannelStockDBOutput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChannelStock", ctx, input, output)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChannelStock indicates an expected call of DeleteChannelStock.
func (mr *MockRepositoryMockRecorder) DeleteChannelStock(ctx, input, output interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChannelStock", reflect.TypeOf((*MockRepository)(nil).DeleteChannelStock), ctx, input, output)
}

// DeleteStock mocks base method.
func (m *MockRepository) DeleteStock(ctx context.Context, input repository.DeleteStockDBInput, output *repository.DeleteStockDBOutput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStock", ctx, input, output)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStock indicates an expected call of DeleteStock.
func (mr *MockRepositoryMockRecorder) DeleteStock(ctx, input, output interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStock", reflect.TypeOf((*MockRepository)(nil).DeleteStock), ctx, input, output)
}

// GetStocks mocks base method.
func (m *MockRepository) GetStocks(ctx context.Context, input repository.GetStocksDBInput, output *repository.GetStocksDBOutput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStocks", ctx, input, output)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetStocks indicates an expected call of GetStocks.
func (mr *MockRepositoryMockRecorder) GetStocks(ctx, input, output interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStocks", reflect.TypeOf((*MockRepository)(nil).GetStocks), ctx, input, output)
}

// UpdateChannelStocks mocks base method.
func (m *MockRepository) UpdateChannelStocks(ctx context.Context, input repository.UpdateChannelStocksDBInput, output *repository.UpdateChannelStocksDBOutput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChannelStocks", ctx, input, output)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateChannelStocks indicates an expected call of UpdateChannelStocks.
func (mr *MockRepositoryMockRecorder) UpdateChannelStocks(ctx, input, output interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChannelStocks", reflect.TypeOf((*MockRepository)(nil).UpdateChannelStocks), ctx, input, output)
}

// UpsertStock mocks base method.
func (m *MockRepository) UpsertStock(ctx context.Context, input repository.UpsertStockDBInput, output *repository.UpsertStockDBOutput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertStock", ctx, input, output)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertStock indicates an expected call of UpsertStock.
func (mr *MockRepositoryMockRecorder) UpsertStock(ctx, input, output interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertStock", reflect.TypeOf((*MockRepository)(nil).UpsertStock), ctx, input, output)
}
