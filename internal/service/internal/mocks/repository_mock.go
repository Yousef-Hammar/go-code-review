// Code generated by mockery v2.40.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/Yousef-Hammar/go-code-review/coupon_service/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// FindByCode provides a mock function with given fields: _a0, _a1
func (_m *Repository) FindByCode(_a0 context.Context, _a1 string) (*domain.Coupon, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for FindByCode")
	}

	var r0 *domain.Coupon
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Coupon, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Coupon); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Coupon)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_FindByCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByCode'
type Repository_FindByCode_Call struct {
	*mock.Call
}

// FindByCode is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *Repository_Expecter) FindByCode(_a0 interface{}, _a1 interface{}) *Repository_FindByCode_Call {
	return &Repository_FindByCode_Call{Call: _e.mock.On("FindByCode", _a0, _a1)}
}

func (_c *Repository_FindByCode_Call) Run(run func(_a0 context.Context, _a1 string)) *Repository_FindByCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_FindByCode_Call) Return(_a0 *domain.Coupon, _a1 error) *Repository_FindByCode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_FindByCode_Call) RunAndReturn(run func(context.Context, string) (*domain.Coupon, error)) *Repository_FindByCode_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *Repository) Save(_a0 context.Context, _a1 domain.Coupon) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Coupon) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type Repository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 domain.Coupon
func (_e *Repository_Expecter) Save(_a0 interface{}, _a1 interface{}) *Repository_Save_Call {
	return &Repository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *Repository_Save_Call) Run(run func(_a0 context.Context, _a1 domain.Coupon)) *Repository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Coupon))
	})
	return _c
}

func (_c *Repository_Save_Call) Return(_a0 error) *Repository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_Save_Call) RunAndReturn(run func(context.Context, domain.Coupon) error) *Repository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
