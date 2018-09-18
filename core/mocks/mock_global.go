// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/iost-official/Go-IOS-Protocol/core/global (interfaces: BaseVariable)

// Package core_mock is a generated GoMock package.
package core_mock

import (
	gomock "github.com/golang/mock/gomock"
	common "github.com/iost-official/Go-IOS-Protocol/common"
	block "github.com/iost-official/Go-IOS-Protocol/core/block"
	global "github.com/iost-official/Go-IOS-Protocol/core/global"
	db "github.com/iost-official/Go-IOS-Protocol/db"
	reflect "reflect"
)

// MockBaseVariable is a mock of BaseVariable interface
type MockBaseVariable struct {
	ctrl     *gomock.Controller
	recorder *MockBaseVariableMockRecorder
}

// MockBaseVariableMockRecorder is the mock recorder for MockBaseVariable
type MockBaseVariableMockRecorder struct {
	mock *MockBaseVariable
}

// NewMockBaseVariable creates a new mock instance
func NewMockBaseVariable(ctrl *gomock.Controller) *MockBaseVariable {
	mock := &MockBaseVariable{ctrl: ctrl}
	mock.recorder = &MockBaseVariableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBaseVariable) EXPECT() *MockBaseVariableMockRecorder {
	return m.recorder
}

// BlockChain mocks base method
func (m *MockBaseVariable) BlockChain() block.Chain {
	ret := m.ctrl.Call(m, "BlockChain")
	ret0, _ := ret[0].(block.Chain)
	return ret0
}

// BlockChain indicates an expected call of BlockChain
func (mr *MockBaseVariableMockRecorder) BlockChain() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockChain", reflect.TypeOf((*MockBaseVariable)(nil).BlockChain))
}

// Config mocks base method
func (m *MockBaseVariable) Config() *common.Config {
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*common.Config)
	return ret0
}

// Config indicates an expected call of Config
func (mr *MockBaseVariableMockRecorder) Config() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockBaseVariable)(nil).Config))
}

// Mode mocks base method
func (m *MockBaseVariable) Mode() global.TMode {
	ret := m.ctrl.Call(m, "Mode")
	ret0, _ := ret[0].(global.TMode)
	return ret0
}

// Mode indicates an expected call of Mode
func (mr *MockBaseVariableMockRecorder) Mode() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mode", reflect.TypeOf((*MockBaseVariable)(nil).Mode))
}

// SetMode mocks base method
func (m *MockBaseVariable) SetMode(arg0 global.TMode) {
	m.ctrl.Call(m, "SetMode", arg0)
}

// SetMode indicates an expected call of SetMode
func (mr *MockBaseVariableMockRecorder) SetMode(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMode", reflect.TypeOf((*MockBaseVariable)(nil).SetMode), arg0)
}

// StateDB mocks base method
func (m *MockBaseVariable) StateDB() db.MVCCDB {
	ret := m.ctrl.Call(m, "StateDB")
	ret0, _ := ret[0].(db.MVCCDB)
	return ret0
}

// StateDB indicates an expected call of StateDB
func (mr *MockBaseVariableMockRecorder) StateDB() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateDB", reflect.TypeOf((*MockBaseVariable)(nil).StateDB))
}

// TxDB mocks base method
func (m *MockBaseVariable) TxDB() global.TxDB {
	ret := m.ctrl.Call(m, "TxDB")
	ret0, _ := ret[0].(global.TxDB)
	return ret0
}

// TxDB indicates an expected call of TxDB
func (mr *MockBaseVariableMockRecorder) TxDB() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxDB", reflect.TypeOf((*MockBaseVariable)(nil).TxDB))
}

// WitnessList mocks base method
func (m *MockBaseVariable) WitnessList() []string {
	ret := m.ctrl.Call(m, "WitnessList")
	ret0, _ := ret[0].([]string)
	return ret0
}

// WitnessList indicates an expected call of WitnessList
func (mr *MockBaseVariableMockRecorder) WitnessList() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WitnessList", reflect.TypeOf((*MockBaseVariable)(nil).WitnessList))
}
