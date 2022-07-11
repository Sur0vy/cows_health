package dataprocessor

import (
	"context"
	"errors"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/logger"
	mock "github.com/Sur0vy/cows_health.git/mocks"
)

func TestDataProcessorLite_CalculateHealth(t *testing.T) {
	type args struct {
		data models.MonitoringDataFull
	}
	type want struct {
		data models.Health
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success",
			args: args{
				data: models.MonitoringDataFull{
					Data: []models.MonitoringData{
						{
							PH:          10,
							Temperature: 20,
							Movement:    30,
						},
						{
							PH:          16,
							Temperature: 22,
							Movement:    41,
						},
					},
					AvgPH:       13,
					AvgTemp:     21,
					AvgMovement: 35.5,
				},
			},
			want: want{
				data: models.Health{
					Estrus: false,
					Ill:    "",
				},
			},
		},
		{
			name: "success",
			args: args{
				data: models.MonitoringDataFull{
					Data: []models.MonitoringData{
						{
							PH:          10,
							Temperature: 20,
							Movement:    20,
						},
						{
							PH:          16,
							Temperature: 30,
							Movement:    51,
						},
					},
					AvgPH:       13,
					AvgTemp:     25,
					AvgMovement: 35.5,
				},
			},
			want: want{
				data: models.Health{
					Estrus: true,
					Ill:    "Инфекционное заболевание",
				},
			},
		},
		{
			name: "success",
			args: args{
				data: models.MonitoringDataFull{
					Data: []models.MonitoringData{
						{
							PH:          5,
							Temperature: 40.5,
							Movement:    30,
						},
						{
							PH:          4,
							Temperature: 40.9,
							Movement:    41,
						},
					},
					AvgPH:       4.5,
					AvgTemp:     20.7,
					AvgMovement: 35.5,
				},
			},
			want: want{
				data: models.Health{
					Estrus: true,
					Ill:    "Нервное заболевание",
				},
			},
		},
	}

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	})

	repoMD := &mock.MonitoringDataStorage{}
	repoCow := &mock.CowStorage{}

	md := NewProcessor(repoMD, repoCow, logger.New())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.data.UpdatedAt = time.Now()
			res := md.CalculateHealth(tt.args.data)
			assert.Equal(t, tt.want.data, res)
		})
	}
}

func TestDataProcessorLite_GetHealthData(t *testing.T) {
	type args struct {
		id       int
		interval int
		data     []models.MonitoringData
		err      error
	}
	type want struct {
		data    models.MonitoringDataFull
		err     error
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success",
			args: args{
				id:       1,
				interval: 10,
				data: []models.MonitoringData{
					{
						PH:          10,
						Temperature: 20,
						Movement:    30,
					},
					{
						PH:          16,
						Temperature: 22,
						Movement:    41,
					},
				},
				err: nil,
			},
			want: want{
				data: models.MonitoringDataFull{
					Data: []models.MonitoringData{
						{
							PH:          10,
							Temperature: 20,
							Movement:    30,
						},
						{
							PH:          16,
							Temperature: 22,
							Movement:    41,
						},
					},
					AvgPH:       13,
					AvgTemp:     21,
					AvgMovement: 35.5,
				},
				err:     nil,
				wantErr: false,
			},
		},
		{
			name: "err",
			args: args{
				id:       1,
				interval: 10,
				err:      errors.New("some error"),
			},
			want: want{
				wantErr: true,
			},
		},
	}

	repoMD := &mock.MonitoringDataStorage{}
	repoCow := &mock.CowStorage{}

	md := NewProcessor(repoMD, repoCow, logger.New())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repoMD.On("Get", context.Background(), tt.args.id, tt.args.interval).
				Return(tt.args.data, tt.args.err).
				Once()

			res, err := md.GetHealthData(context.Background(), tt.args.id)
			if !tt.want.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want.data, res)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestDataProcessorLite_Save(t *testing.T) {
	type args struct {
		data models.MonitoringData
		err  error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "right",
			args: args{
				data: models.MonitoringData{
					ID:          1,
					BolusNum:    1,
					CowID:       1,
					AddedAt:     time.Time{},
					PH:          5,
					Temperature: 35,
					Movement:    20,
					Charge:      99,
				},
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				data: models.MonitoringData{
					ID:          1,
					BolusNum:    1,
					CowID:       1,
					AddedAt:     time.Time{},
					PH:          5,
					Temperature: 35,
					Movement:    20,
					Charge:      99,
				},
				err: errors.New("some error"),
			},
			wantErr: true,
		},
	}

	repoMD := &mock.MonitoringDataStorage{}
	repoCow := &mock.CowStorage{}

	md := NewProcessor(repoMD, repoCow, logger.New())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repoMD.On("Add", context.Background(), tt.args.data).
				Return(tt.args.err).
				Once()

			err := md.Save(context.Background(), tt.args.data)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
