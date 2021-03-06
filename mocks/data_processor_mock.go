// Code generated by mockery v2.13.1. DO NOT EDIT.

package internalMock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/Sur0vy/cows_health.git/internal/models"

	sync "sync"
)

// DataProcessor is an autogenerated mock type for the DataProcessor type
type DataProcessor struct {
	mock.Mock
}

// CalculateHealth provides a mock function with given fields: mds
func (_m *DataProcessor) CalculateHealth(mds models.MonitoringDataFull) models.Health {
	ret := _m.Called(mds)

	var r0 models.Health
	if rf, ok := ret.Get(0).(func(models.MonitoringDataFull) models.Health); ok {
		r0 = rf(mds)
	} else {
		r0 = ret.Get(0).(models.Health)
	}

	return r0
}

// GetHealthData provides a mock function with given fields: c, cowID
func (_m *DataProcessor) GetHealthData(c context.Context, cowID int) (models.MonitoringDataFull, error) {
	ret := _m.Called(c, cowID)

	var r0 models.MonitoringDataFull
	if rf, ok := ret.Get(0).(func(context.Context, int) models.MonitoringDataFull); ok {
		r0 = rf(c, cowID)
	} else {
		r0 = ret.Get(0).(models.MonitoringDataFull)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(c, cowID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Run provides a mock function with given fields: c, data, wg
func (_m *DataProcessor) Run(c context.Context, data models.MonitoringData, wg *sync.WaitGroup) {
	_m.Called(c, data, wg)
}

// Save provides a mock function with given fields: c, data
func (_m *DataProcessor) Save(c context.Context, data models.MonitoringData) error {
	ret := _m.Called(c, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.MonitoringData) error); ok {
		r0 = rf(c, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDataProcessor interface {
	mock.TestingT
	Cleanup(func())
}

// NewDataProcessor creates a new instance of DataProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDataProcessor(t mockConstructorTestingTNewDataProcessor) *DataProcessor {
	mock := &DataProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
