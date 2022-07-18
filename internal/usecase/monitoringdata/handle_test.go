package monitoringdata

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/logger"
	mock "github.com/Sur0vy/cows_health.git/mocks"
)

func TestHandler_Add(t *testing.T) {
	type args struct {
		useMocks bool
		data     []models.MonitoringData
		badBody  string
		err      error
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "bad",
			args: args{
				useMocks: false,
				badBody:  "bad body",
				err:      nil,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "success",
			args: args{
				useMocks: true,
				err:      nil,
			},
			want: want{
				code: http.StatusAccepted,
			},
		},
	}

	repoMD := &mock.MonitoringDataStorage{}
	repoDP := &mock.DataProcessor{}

	mdh := NewMonitoringDataHandler(repoMD, repoDP, logger.New())
	router := echo.New()
	router.POST("api/data", mdh.Add)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.args.useMocks {
				repoMD.On("Add", context.Background(), tt.args.data).
					Return(tt.args.err).
					Once()
				body, _ = json.Marshal(tt.args.data)
			} else {
				body = []byte(tt.args.badBody)
			}
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/data", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
		})
	}
}
