// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SpellerValidator is an autogenerated mock type for the SpellerValidator type
type SpellerValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: texts
func (_m *SpellerValidator) Validate(texts []string) error {
	ret := _m.Called(texts)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(texts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewSpellerValidator interface {
	mock.TestingT
	Cleanup(func())
}

// NewSpellerValidator creates a new instance of SpellerValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSpellerValidator(t mockConstructorTestingTNewSpellerValidator) *SpellerValidator {
	mock := &SpellerValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
