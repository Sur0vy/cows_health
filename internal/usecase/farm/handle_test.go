package farm

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Sur0vy/cows_health.git/errors"
	"github.com/Sur0vy/cows_health.git/logger"
	"github.com/Sur0vy/cows_health.git/mocks"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

func TestHandler_Add(t *testing.T) {
	type args struct {
		useMocks bool
		useUser  bool
		farm     models.Farm
		userID   int
		body     string
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
			name: "add success",
			args: args{
				useMocks: true,
				useUser:  true,
				farm: models.Farm{
					Name:    "Farm_1",
					Address: "Farm_1_address",
					UserID:  1,
				},
				userID: 1,
				body:   "{\"name\": \"Farm_1\",\"address\": \"Farm_1_address\"}",
				err:    nil,
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "no user",
			args: args{
				useMocks: false,
				useUser:  false,
				farm: models.Farm{
					Name:    "Farm_1",
					Address: "Farm_1_address",
					UserID:  1,
				},
				userID: 1,
				body:   "{\"name\": \"Farm_1\",\"address\": \"Farm_1_address\"}",
				err:    nil,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "bad body",
			args: args{
				useMocks: false,
				useUser:  false,
				farm: models.Farm{
					Name:    "Farm_1",
					Address: "Farm_1_address",
					UserID:  1,
				},
				userID: 1,
				body:   "{\"namFarm_1\",\"address\": \"Farm_1",
				err:    nil,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "the same farm",
			args: args{
				useMocks: true,
				useUser:  true,
				farm: models.Farm{
					Name:    "Farm_1",
					Address: "Farm_1_address",
					UserID:  1,
				},
				userID: 1,
				body:   "{\"name\": \"Farm_1\",\"address\": \"Farm_1_address\"}",
				err:    errors.ErrExist,
			},
			want: want{
				code: http.StatusConflict,
			},
		},
		{
			name: "db request corrupted farm",
			args: args{
				useMocks: true,
				useUser:  true,
				farm: models.Farm{
					Name:    "Farm_1",
					Address: "Farm_1_address",
					UserID:  1,
				},
				userID: 1,
				body:   "{\"name\": \"Farm_1\",\"address\": \"Farm_1_address\"}",
				err:    errors.ErrEmpty,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	repo := &storageMock.FarmStorage{}

	fh := NewFarmHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/farm", fh.Add)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					if tt.args.useUser {
						c.Set("UserID", tt.args.userID)
					}
					return next(c)
				}
			})

			if tt.args.useMocks {
				repo.On("Add", context.Background(), tt.args.farm).
					Return(tt.args.err).
					Once()
			}
			body := bytes.NewBuffer([]byte(tt.args.body))
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/farm", body)
			if err != nil {
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						panic(err)
					}
				}(req.Body)
			}

			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type args struct {
		useMocks bool
		farmID   string
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
			name: "delete success",
			args: args{
				useMocks: true,
				farmID:   "1",
				err:      nil,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "bad request",
			args: args{
				useMocks: false,
				farmID:   "bad",
				err:      nil,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "no farm",
			args: args{
				useMocks: true,
				farmID:   "1",
				err:      errors.ErrEmpty,
			},
			want: want{
				code: http.StatusConflict,
			},
		},
		{
			name: "db request corrupted",
			args: args{
				useMocks: true,
				farmID:   "1",
				err:      errors.ErrExist,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	repo := &storageMock.FarmStorage{}

	fh := NewFarmHandler(repo, logger.New())
	router := echo.New()
	router.DELETE("/api/farm/:id", fh.Delete)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.useMocks {
				val, _ := strconv.Atoi(tt.args.farmID)
				repo.On("Delete", context.Background(), val).
					Return(tt.args.err).
					Once()
			}
			recorder := httptest.NewRecorder()
			path := fmt.Sprintf("/api/farm/%s", tt.args.farmID)
			req, err := http.NewRequest("DELETE", path, nil)
			if err != nil {
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						panic(err)
					}
				}(req.Body)
			}

			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
		})
	}
}

func TestHandler_Get(t *testing.T) {
	type args struct {
		useMocks bool
		userID   int
		farms    []models.Farm
		err      error
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "get success",
			args: args{
				useMocks: true,
				userID:   1,
				farms: []models.Farm{
					{
						ID:      1,
						Name:    "Farm_1",
						Address: "Address_1",
					},
					{
						ID:      2,
						Name:    "Farm_2",
						Address: "Address_2",
					},
				},
				err: nil,
			},
			want: want{
				code: http.StatusOK,
				body: "[{\"farm_id\":1,\"name\":\"Farm_1\",\"address\":\"Address_1\"},{\"farm_id\":2,\"name\":\"Farm_2\",\"address\":\"Address_2\"}]\n",
			},
		},
		{
			name: "no user",
			args: args{
				useMocks: false,
				userID:   1,
				err:      nil,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "no farms",
			args: args{
				useMocks: true,
				userID:   1,
				farms: []models.Farm{
					{
						ID:      1,
						Name:    "Farm_1",
						Address: "Address_1",
					},
					{
						ID:      2,
						Name:    "Farm_2",
						Address: "Address_2",
					},
				},
				err: errors.ErrEmpty,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "db request corrupted",
			args: args{
				useMocks: true,
				userID:   1,
				farms: []models.Farm{
					{
						ID:      1,
						Name:    "Farm_1",
						Address: "Address_1",
						UserID:  1,
					},
					{
						ID:      2,
						Name:    "Farm_2",
						Address: "Address_2",
						UserID:  1,
					},
				},
				err: errors.ErrExist,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	repo := &storageMock.FarmStorage{}

	fh := NewFarmHandler(repo, logger.New())
	router := echo.New()
	router.GET("/api/farm", fh.Get)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					if tt.args.useMocks {
						c.Set("UserID", tt.args.userID)
					}
					return next(c)
				}
			})
			if tt.args.useMocks {
				repo.On("Get", context.Background(), tt.args.userID).
					Return(tt.args.farms, tt.args.err).
					Once()
			}
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/farm", nil)
			if err != nil {
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						panic(err)
					}
				}(req.Body)
			}

			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
			if tt.want.code == http.StatusOK {
				input, _ := ioutil.ReadAll(recorder.Result().Body)
				assert.Equal(t, tt.want.body, string(input))
			}
		})
	}
}
