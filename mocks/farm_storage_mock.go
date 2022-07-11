// Code generated by mockery v2.13.1. DO NOT EDIT.

package internalMock

import (
	context "context"

	models "github.com/Sur0vy/cows_health.git/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// FarmStorage is an autogenerated mock type for the FarmStorage type
type FarmStorage struct {
	mock.Mock
}

// Add provides a mock function with given fields: с, farm
func (_m *FarmStorage) Add(с context.Context, farm models.Farm) error {
	ret := _m.Called(с, farm)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Farm) error); ok {
		r0 = rf(с, farm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: с, farmID
func (_m *FarmStorage) Delete(с context.Context, farmID int) error {
	ret := _m.Called(с, farmID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(с, farmID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: с, userID
func (_m *FarmStorage) Get(с context.Context, userID int) ([]models.Farm, error) {
	ret := _m.Called(с, userID)

	var r0 []models.Farm
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.Farm); ok {
		r0 = rf(с, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Farm)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(с, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFarmStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewFarmStorage creates a new instance of FarmStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFarmStorage(t mockConstructorTestingTNewFarmStorage) *FarmStorage {
	mock := &FarmStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
