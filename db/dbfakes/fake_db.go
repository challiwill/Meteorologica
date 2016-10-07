// This file was generated by counterfeiter
package dbfakes

import (
	"database/sql"
	"sync"

	"github.com/challiwill/meteorologica/db"
)

type FakeDB struct {
	ExecStub        func(string, ...interface{}) (sql.Result, error)
	execMutex       sync.RWMutex
	execArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	execReturns struct {
		result1 sql.Result
		result2 error
	}
	QueryRowStub        func(string, ...interface{}) *sql.Row
	queryRowMutex       sync.RWMutex
	queryRowArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	queryRowReturns struct {
		result1 *sql.Row
	}
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns     struct {
		result1 error
	}
	PingStub        func() error
	pingMutex       sync.RWMutex
	pingArgsForCall []struct{}
	pingReturns     struct {
		result1 error
	}
	BeginStub        func() (*sql.Tx, error)
	beginMutex       sync.RWMutex
	beginArgsForCall []struct{}
	beginReturns     struct {
		result1 *sql.Tx
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDB) Exec(arg1 string, arg2 ...interface{}) (sql.Result, error) {
	fake.execMutex.Lock()
	fake.execArgsForCall = append(fake.execArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	fake.recordInvocation("Exec", []interface{}{arg1, arg2})
	fake.execMutex.Unlock()
	if fake.ExecStub != nil {
		return fake.ExecStub(arg1, arg2...)
	} else {
		return fake.execReturns.result1, fake.execReturns.result2
	}
}

func (fake *FakeDB) ExecCallCount() int {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return len(fake.execArgsForCall)
}

func (fake *FakeDB) ExecArgsForCall(i int) (string, []interface{}) {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return fake.execArgsForCall[i].arg1, fake.execArgsForCall[i].arg2
}

func (fake *FakeDB) ExecReturns(result1 sql.Result, result2 error) {
	fake.ExecStub = nil
	fake.execReturns = struct {
		result1 sql.Result
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) QueryRow(arg1 string, arg2 ...interface{}) *sql.Row {
	fake.queryRowMutex.Lock()
	fake.queryRowArgsForCall = append(fake.queryRowArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	fake.recordInvocation("QueryRow", []interface{}{arg1, arg2})
	fake.queryRowMutex.Unlock()
	if fake.QueryRowStub != nil {
		return fake.QueryRowStub(arg1, arg2...)
	} else {
		return fake.queryRowReturns.result1
	}
}

func (fake *FakeDB) QueryRowCallCount() int {
	fake.queryRowMutex.RLock()
	defer fake.queryRowMutex.RUnlock()
	return len(fake.queryRowArgsForCall)
}

func (fake *FakeDB) QueryRowArgsForCall(i int) (string, []interface{}) {
	fake.queryRowMutex.RLock()
	defer fake.queryRowMutex.RUnlock()
	return fake.queryRowArgsForCall[i].arg1, fake.queryRowArgsForCall[i].arg2
}

func (fake *FakeDB) QueryRowReturns(result1 *sql.Row) {
	fake.QueryRowStub = nil
	fake.queryRowReturns = struct {
		result1 *sql.Row
	}{result1}
}

func (fake *FakeDB) Close() error {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	} else {
		return fake.closeReturns.result1
	}
}

func (fake *FakeDB) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeDB) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDB) Ping() error {
	fake.pingMutex.Lock()
	fake.pingArgsForCall = append(fake.pingArgsForCall, struct{}{})
	fake.recordInvocation("Ping", []interface{}{})
	fake.pingMutex.Unlock()
	if fake.PingStub != nil {
		return fake.PingStub()
	} else {
		return fake.pingReturns.result1
	}
}

func (fake *FakeDB) PingCallCount() int {
	fake.pingMutex.RLock()
	defer fake.pingMutex.RUnlock()
	return len(fake.pingArgsForCall)
}

func (fake *FakeDB) PingReturns(result1 error) {
	fake.PingStub = nil
	fake.pingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDB) Begin() (*sql.Tx, error) {
	fake.beginMutex.Lock()
	fake.beginArgsForCall = append(fake.beginArgsForCall, struct{}{})
	fake.recordInvocation("Begin", []interface{}{})
	fake.beginMutex.Unlock()
	if fake.BeginStub != nil {
		return fake.BeginStub()
	} else {
		return fake.beginReturns.result1, fake.beginReturns.result2
	}
}

func (fake *FakeDB) BeginCallCount() int {
	fake.beginMutex.RLock()
	defer fake.beginMutex.RUnlock()
	return len(fake.beginArgsForCall)
}

func (fake *FakeDB) BeginReturns(result1 *sql.Tx, result2 error) {
	fake.BeginStub = nil
	fake.beginReturns = struct {
		result1 *sql.Tx
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	fake.queryRowMutex.RLock()
	defer fake.queryRowMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.pingMutex.RLock()
	defer fake.pingMutex.RUnlock()
	fake.beginMutex.RLock()
	defer fake.beginMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeDB) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ db.DB = new(FakeDB)
