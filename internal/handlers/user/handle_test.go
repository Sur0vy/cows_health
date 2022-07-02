package user

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/Sur0vy/cows_health.git/internal/errors"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	storage_mock "github.com/Sur0vy/cows_health.git/internal/mocks"
	"github.com/Sur0vy/cows_health.git/internal/models"
)

func TestHandler_Login(t *testing.T) {
	//то, что замокано
	type args struct {
		user models.User
		body string
		hash string
		err  error
	}
	//ожидание
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
			name: "login success",
			args: args{
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
		{
			name: "bad body",
			args: args{
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body: "uncorrect JSON string",
				hash: "8f9bfe9d1345237cb3b2b205864da075",
				err:  nil,
			},
			want: want{
				code: http.StatusBadRequest,
				hash: "",
			},
		},
	}

	repo := &storage_mock.UserStorage{}

	uh := NewUserHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/user/login", uh.Login)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetHash", context.Background(), tt.args.user).
				Return(tt.args.hash, tt.args.err).
				Once()
			body := bytes.NewBuffer([]byte(tt.args.body))
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/user/login", body)
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
	//то, что замокано
	type args struct {
		body string
	}
	//ожидание
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

	repo := &storage_mock.UserStorage{}

	uh := NewUserHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/user/logout", uh.Logout)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//repo.On("GetHash", context.Background(), tt.args.user).
			//	Return(tt.args.hash, tt.args.err).
			//	Once()
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
	//то, что замокано
	type args struct {
		user models.User
		body string
		hash string
		err  error
	}
	//ожидание
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
			name: "login success",
			args: args{
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
		{
			name: "bad body",
			args: args{
				user: models.User{
					Login:    "User",
					Password: "pa$$word_1",
				},
				body: "uncorrect JSON string",
				hash: "8f9bfe9d1345237cb3b2b205864da075",
				err:  nil,
			},
			want: want{
				code: http.StatusBadRequest,
				hash: "",
			},
		},
	}

	repo := &storage_mock.UserStorage{}

	uh := NewUserHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/user/login", uh.Login)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetHash", context.Background(), tt.args.user).
				Return(tt.args.hash, tt.args.err).
				Once()
			repo.On("Add", context.Background(), tt.args.user).
				Return(tt.args.hash, tt.args.err).
				Once()
			body := bytes.NewBuffer([]byte(tt.args.body))
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/user/login", body)
			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
			if len(recorder.Result().Cookies()) > 0 && len(tt.want.hash) > 0 {
				assert.Equal(t, tt.want.hash, recorder.Result().Cookies()[0].Value)
			}
		})
	}
}
