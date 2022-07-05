package user

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/Sur0vy/cows_health.git/internal/errors"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	storageMock "github.com/Sur0vy/cows_health.git/internal/mocks"
	"github.com/Sur0vy/cows_health.git/internal/models"
)

func TestHandler_Login(t *testing.T) {
	type args struct {
		useMocks bool
		user     models.User
		body     string
		hash     string
		err      error
	}
	type want struct {
		code int
		hash string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "unexpected storage error",
			args: args{
				useMocks: true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body: "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash: "8f9bfe9d1345237cb3b2b205864da075",
				err:  errors.NewExistError(),
			},
			want: want{
				code: http.StatusInternalServerError,
				hash: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
		{
			name: "no user",
			args: args{
				useMocks: true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body: "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash: "8f9bfe9d1345237cb3b2b205864da075",
				err:  errors.NewEmptyError(),
			},
			want: want{
				code: http.StatusUnauthorized,
				hash: "",
			},
		},
		{
			name: "bad body",
			args: args{
				useMocks: false,
				user:     models.User{},
				body:     "uncorrected JSON string",
				hash:     "8f9bfe9d1345237cb3b2b205864da075",
				err:      nil,
			},
			want: want{
				code: http.StatusBadRequest,
				hash: "",
			},
		},
		{
			name: "login success",
			args: args{
				useMocks: true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body: "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash: "8f9bfe9d1345237cb3b2b205864da075",
				err:  nil,
			},
			want: want{
				code: http.StatusOK,
				hash: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
	}

	repo := &storageMock.UserStorage{}

	uh := NewUserHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/user/login", uh.Login)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.useMocks {
				repo.On("GetHash", context.Background(), tt.args.user).
					Return(tt.args.hash, tt.args.err).
					Once()
			}
			body := bytes.NewBuffer([]byte(tt.args.body))
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/user/login", body)
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
			if len(recorder.Result().Cookies()) > 0 && len(tt.want.hash) > 0 {
				assert.Equal(t, tt.want.hash, recorder.Result().Cookies()[0].Value)
			}
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	type args struct {
		body string
	}
	type want struct {
		code int
		hash string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "logout success",
			args: args{
				body: "",
			},
			want: want{
				code: http.StatusOK,
				hash: "",
			},
		},
	}

	repo := &storageMock.UserStorage{}

	uh := NewUserHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/user/logout", uh.Logout)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(tt.args.body))
			req, err := http.NewRequest("POST", "/api/user/logout", body)
			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
			if len(recorder.Result().Cookies()) > 0 && len(tt.want.hash) > 0 {
				assert.Equal(t, tt.want.hash, recorder.Result().Cookies()[0].Value)
			}
		})
	}
}

func TestHandler_Register(t *testing.T) {
	type args struct {
		useMocksGH bool
		useMocksA  bool
		user       models.User
		body       string
		hash       string
		errGH      error
		errA       error
	}
	type want struct {
		code int
		hash string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "add user success",
			args: args{
				useMocksGH: true,
				useMocksA:  true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body:  "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash:  "8f9bfe9d1345237cb3b2b205864da075",
				errGH: nil,
				errA:  nil,
			},
			want: want{
				code: http.StatusOK,
				hash: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
		{
			name: "bad body",
			args: args{
				useMocksGH: false,
				useMocksA:  false,
				user:       models.User{},
				body:       "uncorrected JSON string",
				hash:       "8f9bfe9d1345237cb3b2b205864da075",
				errGH:      nil,
				errA:       nil,
			},
			want: want{
				code: http.StatusBadRequest,
				hash: "",
			},
		},
		{
			name: "user exist",
			args: args{
				useMocksGH: false,
				useMocksA:  true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body:  "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash:  "8f9bfe9d1345237cb3b2b205864da075",
				errGH: nil,
				errA:  errors.NewExistError(),
			},
			want: want{
				code: http.StatusConflict,
				hash: "",
			},
		},
		{
			name: "unexpected error add",
			args: args{
				useMocksGH: false,
				useMocksA:  true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body:  "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash:  "8f9bfe9d1345237cb3b2b205864da075",
				errGH: nil,
				errA:  errors.NewEmptyError(),
			},
			want: want{
				code: http.StatusInternalServerError,
				hash: "",
			},
		},
		{
			name: "get user hash error",
			args: args{
				useMocksGH: true,
				useMocksA:  true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body:  "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash:  "8f9bfe9d1345237cb3b2b205864da075",
				errGH: errors.NewEmptyError(),
				errA:  nil,
			},
			want: want{
				code: http.StatusUnauthorized,
				hash: "",
			},
		},
		{
			name: "unexpected storage error",
			args: args{
				useMocksGH: true,
				useMocksA:  true,
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body:  "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
				hash:  "8f9bfe9d1345237cb3b2b205864da075",
				errGH: errors.NewExistError(),
				errA:  nil,
			},
			want: want{
				code: http.StatusInternalServerError,
				hash: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
	}

	repo := &storageMock.UserStorage{}

	uh := NewUserHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/user/register", uh.Register)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.useMocksGH {
				repo.On("GetHash", context.Background(), tt.args.user).
					Return(tt.args.hash, tt.args.errGH).
					Once()
			}
			if tt.args.useMocksA {
				repo.On("Add", context.Background(), tt.args.user).
					Return(tt.args.errA).
					Once()
			}
			body := bytes.NewBuffer([]byte(tt.args.body))
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/user/register", body)
			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
			if len(recorder.Result().Cookies()) > 0 && len(tt.want.hash) > 0 {
				assert.Equal(t, tt.want.hash, recorder.Result().Cookies()[0].Value)
			}
		})
	}
}
