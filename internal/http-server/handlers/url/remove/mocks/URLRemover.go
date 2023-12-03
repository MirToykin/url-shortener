// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// URLRemover is an autogenerated mock type for the URLRemover type
type URLRemover struct {
	mock.Mock
}

// DeleteUrl provides a mock function with given fields: alias
func (_m *URLRemover) DeleteUrl(alias string) error {
	ret := _m.Called(alias)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewURLRemover interface {
	mock.TestingT
	Cleanup(func())
}

// NewURLRemover creates a new instance of URLRemover. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewURLRemover(t mockConstructorTestingTNewURLRemover) *URLRemover {
	mock := &URLRemover{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
