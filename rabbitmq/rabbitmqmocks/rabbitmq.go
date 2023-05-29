// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/getAlby/lndhub.go/rabbitmq (interfaces: LndHubService)

// Package rabbitmqmocks is a generated GoMock package.
package rabbitmqmocks

import (
	context "context"
	reflect "reflect"

	models "github.com/getAlby/lndhub.go/db/models"
	gomock "github.com/golang/mock/gomock"
)

// MockLndHubService is a mock of LndHubService interface.
type MockLndHubService struct {
	ctrl     *gomock.Controller
	recorder *MockLndHubServiceMockRecorder
}

// MockLndHubServiceMockRecorder is the mock recorder for MockLndHubService.
type MockLndHubServiceMockRecorder struct {
	mock *MockLndHubService
}

// NewMockLndHubService creates a new mock instance.
func NewMockLndHubService(ctrl *gomock.Controller) *MockLndHubService {
	mock := &MockLndHubService{ctrl: ctrl}
	mock.recorder = &MockLndHubServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLndHubService) EXPECT() *MockLndHubServiceMockRecorder {
	return m.recorder
}

// GetAllPendingPayments mocks base method.
func (m *MockLndHubService) GetAllPendingPayments(arg0 context.Context) ([]models.Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPendingPayments", arg0)
	ret0, _ := ret[0].([]models.Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPendingPayments indicates an expected call of GetAllPendingPayments.
func (mr *MockLndHubServiceMockRecorder) GetAllPendingPayments(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPendingPayments", reflect.TypeOf((*MockLndHubService)(nil).GetAllPendingPayments), arg0)
}

// GetTransactionEntryByInvoiceId mocks base method.
func (m *MockLndHubService) GetTransactionEntryByInvoiceId(arg0 context.Context, arg1 int64) (models.TransactionEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionEntryByInvoiceId", arg0, arg1)
	ret0, _ := ret[0].(models.TransactionEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionEntryByInvoiceId indicates an expected call of GetTransactionEntryByInvoiceId.
func (mr *MockLndHubServiceMockRecorder) GetTransactionEntryByInvoiceId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionEntryByInvoiceId", reflect.TypeOf((*MockLndHubService)(nil).GetTransactionEntryByInvoiceId), arg0, arg1)
}

// HandleFailedPayment mocks base method.
func (m *MockLndHubService) HandleFailedPayment(arg0 context.Context, arg1 *models.Invoice, arg2 models.TransactionEntry, arg3 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleFailedPayment", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleFailedPayment indicates an expected call of HandleFailedPayment.
func (mr *MockLndHubServiceMockRecorder) HandleFailedPayment(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleFailedPayment", reflect.TypeOf((*MockLndHubService)(nil).HandleFailedPayment), arg0, arg1, arg2, arg3)
}

// HandleSuccessfulPayment mocks base method.
func (m *MockLndHubService) HandleSuccessfulPayment(arg0 context.Context, arg1 *models.Invoice, arg2 models.TransactionEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleSuccessfulPayment", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleSuccessfulPayment indicates an expected call of HandleSuccessfulPayment.
func (mr *MockLndHubServiceMockRecorder) HandleSuccessfulPayment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleSuccessfulPayment", reflect.TypeOf((*MockLndHubService)(nil).HandleSuccessfulPayment), arg0, arg1, arg2)
}
