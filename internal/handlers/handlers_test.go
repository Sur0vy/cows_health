package handlers

import (
	"bytes"
	"context"
	"github.com/Sur0vy/cows_health.git/internal/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func TestBaseHandler_AddFarm(t *testing.T) {
	type want struct {
		cookie string
		code   int
		body   string
	}
	testData := []struct {
		body   string
		cookie string
	}{
		{
			body:   "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			cookie: "8f9bfe9d1345237cb3b2b205864da075",
		},
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "New farm",
			args: "{\"name\": \"Farm\",\"address\": \"Moscow\"}",
			want: want{
				code:   http.StatusCreated,
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "",
			},
		},
		{
			name: "Duplicate farm",
			args: "{\"name\": \"Farm\",\"address\": \"Moscow\"}",
			want: want{
				code:   http.StatusConflict,
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "",
			},
		},
		{
			name: "Bad request",
			args: "\"Farm\",\"address\": \"Moscow\"}",
			want: want{
				code:   http.StatusBadRequest,
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "",
			},
		},
		{
			name: "Bad cookie",
			args: "{\"name\": \"Farm\",\"address\": \"Moscow\"}",
			want: want{
				code:   http.StatusUnauthorized,
				cookie: "8f9bfe932423423b2b205864da075",
				body:   "{\"Message\":\"Unauthorized\"}",
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.POST("/farms", handler.AddFarm)
	router.NoRoute(handler.ResponseBadRequest)
	w := httptest.NewRecorder()

	router.POST("/api/user/register", handler.Register)
	for _, tt := range testData {
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.want.cookie, HttpOnly: true}

			body := bytes.NewBuffer([]byte(tt.args))
			req, err := http.NewRequest("POST", "/farms", body)
			req.AddCookie(cookie)

			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			if w.Code == http.StatusUnauthorized {
				input, _ := ioutil.ReadAll(w.Body)
				assert.Equal(t, tt.want.body, string(input))
			}
		})
	}
}

func TestBaseHandler_DelFarm(t *testing.T) {
	assert.NotNil(t, nil)
}

func TestBaseHandler_GetFarms(t *testing.T) {
	assert.NotNil(t, nil)
}

func TestBaseHandler_Login(t *testing.T) {
	type want struct {
		cookie string
		code   int
		err    error
	}
	testData := []struct {
		body string
	}{
		{
			body: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
		},
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "New user",
			args: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			want: want{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				code:   http.StatusOK,
				err:    nil,
			},
		},
		{
			name: "Bad body",
			args: "word\": \"pa$$1\"}",
			want: want{
				cookie: "",
				code:   http.StatusBadRequest,
				err:    nil,
			},
		},
		{
			name: "Wrong user name",
			args: "{\"login\": \"User2\", \"password\": \"pa$$word_1\"}",
			want: want{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				code:   http.StatusUnauthorized,
				err:    nil,
			},
		},
		{
			name: "Wrong password",
			args: "{\"login\": \"User\", \"password\": \"ewrwerwer\"}",
			want: want{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				code:   http.StatusUnauthorized,
				err:    nil,
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.POST("/api/user/login", handler.Login)
	router.NoRoute(handler.ResponseBadRequest)
	w := httptest.NewRecorder()

	router.POST("/api/user/register", handler.Register)
	for _, tt := range testData {
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(tt.args))
			req, err := http.NewRequest("POST", "/api/user/login", body)
			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestBaseHandler_Logout(t *testing.T) {
	type want struct {
		cookie string
		code   int
		err    error
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Logout",
			want: want{
				cookie: "",
				code:   http.StatusOK,
				err:    nil,
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.POST("/api/user/logout", handler.Logout)
	router.NoRoute(handler.ResponseBadRequest)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			//req, err := http.NewRequest("POST", "/api/user/logout", bytes.NewBuffer([]byte("")))
			req, err := http.NewRequest("POST", "/api/user/logout", nil)
			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.cookie, w.Result().Cookies()[0].Value)
		})
	}
}

func TestBaseHandler_Register(t *testing.T) {
	type want struct {
		cookie string
		code   int
		err    error
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "New user",
			args: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			want: want{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				code:   http.StatusOK,
				err:    nil,
			},
		},
		{
			name: "Duplicate user",
			args: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			want: want{
				cookie: "",
				code:   http.StatusConflict,
				err:    nil,
			},
		},
		{
			name: "Bad login",
			args: "{\"lgin\": \"User\", \"password\": \"pa$$rd_1\"}",
			want: want{
				cookie: "",
				code:   http.StatusBadRequest,
				err:    nil,
			},
		},
		{
			name: "Bad body",
			args: "word\": \"pa$$1\"}",
			want: want{
				cookie: "",
				code:   http.StatusBadRequest,
				err:    nil,
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.POST("/api/user/register", handler.Register)
	router.NoRoute(handler.ResponseBadRequest)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(tt.args))
			req, err := http.NewRequest("POST", "/api/user/register", body)
			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			if w.Code == http.StatusOK {
				assert.Equal(t, tt.want.cookie, w.Result().Cookies()[0].Value)
			}
		})
	}
}

func Test_getIDFromJSON(t *testing.T) {
	assert.NotNil(t, nil)
}
