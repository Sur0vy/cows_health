package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Sur0vy/cows_health.git/internal/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func TestBaseHandler_AddFarm(t *testing.T) {
	testData := []struct {
		body   string
		cookie string
	}{
		{
			body:   "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			cookie: "8f9bfe9d1345237cb3b2b205864da075",
		},
	}
	type want struct {
		code int
		body string
	}
	type args struct {
		cookie string
		body   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "New farm",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "{\"name\": \"Farm\",\"address\": \"Moscow\"}",
			},
			want: want{
				code: http.StatusCreated,
				body: "",
			},
		},
		{
			name: "Duplicate farm",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "{\"name\": \"Farm\",\"address\": \"Moscow\"}",
			},
			want: want{
				code: http.StatusConflict,
				body: "",
			},
		},
		{
			name: "Bad request",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "\"Farm\",\"address\": \"Moscow\"}",
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
		{
			name: "Bad cookie",
			args: args{
				cookie: "8f9bfe932423423b2b205864da075",
				body:   "{\"name\": \"Farm\",\"address\": \"Moscow\"}",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "{\"Message\":\"Unauthorized\"}",
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.POST("/api/farms", handler.AddFarm)
	router.NoRoute(handler.ResponseBadRequest)

	router.POST("/api/user/register", handler.Register)
	for _, tt := range testData {
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.args.cookie, HttpOnly: true}

			body := bytes.NewBuffer([]byte(tt.args.body))
			req, err := http.NewRequest("POST", "/api/farms", body)
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
	type user struct {
		body string
	}
	type farm struct {
		body   string
		cookie string
	}
	reliableData := struct {
		users []user
		farms []farm
	}{
		users: []user{
			{
				body: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			},
			{
				body: "{\"login\": \"User2\", \"password\": \"pa$$word_2\"}",
			},
		},
		farms: []farm{
			{
				body:   "{\"name\": \"Farm1\",\"address\": \"Moscow\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"Farm2\",\"address\": \"Omsk\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
	}

	type want struct {
		code int
		body string
	}
	type args struct {
		cookie string
		ID     string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Del farm user1",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				ID:     "1",
			},
			want: want{
				code: http.StatusOK,
				body: "",
			},
		},
		{
			name: "Del not exists farm",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				ID:     "3",
			},
			want: want{
				code: http.StatusConflict,
				body: "",
			},
		},
		{
			name: "No user",
			args: args{
				cookie: "a09bccf2b294353452b34dc0e08d8b582a",
				ID:     "5",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "{\"Message\":\"Unauthorized\"}",
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.DELETE("/api/farms/:id", handler.DelFarm)
	router.NoRoute(handler.ResponseBadRequest)

	//пропишем пользователя
	router.POST("/api/user/register", handler.Register)
	for _, tt := range reliableData.users {
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}
	//пропишем фермы
	router.POST("/api/farms", handler.AddFarm)
	for _, tt := range reliableData.farms {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{Name: config.Cookie, Value: tt.cookie, HttpOnly: true}
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/farms", body)
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
	}

	//запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.args.cookie, HttpOnly: true}
			URL := "/api/farms/" + tt.args.ID
			req, err := http.NewRequest("DELETE", URL, nil)
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

func TestBaseHandler_GetFarms(t *testing.T) {
	type user struct {
		body string
	}
	type farm struct {
		body   string
		cookie string
	}
	reliableData := struct {
		users []user
		farms []farm
	}{
		users: []user{
			{
				body: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			},
			{
				body: "{\"login\": \"User2\", \"password\": \"pa$$word_2\"}",
			},
		},
		farms: []farm{
			{
				body:   "{\"name\": \"Farm1\",\"address\": \"Moscow\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"Farm2\",\"address\": \"Omsk\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
	}

	type want struct {
		code int
		body string
	}
	type args struct {
		cookie string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get farm user1",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			want: want{
				code: http.StatusOK,
				body: "[{\"farm_id\":1,\"name\":\"Farm1\",\"address\":\"Moscow\"}," +
					"{\"farm_id\":2,\"name\":\"Farm2\",\"address\":\"Omsk\"}]",
			},
		},
		{
			name: "Get farm user2",
			args: args{
				cookie: "a09bccf2b2963982b34dc0e08d8b582a",
			},
			want: want{
				code: http.StatusNoContent,
				body: "",
			},
		},
		{
			name: "No user",
			args: args{
				cookie: "a09bccf2b294353452b34dc0e08d8b582a",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "{\"Message\":\"Unauthorized\"}",
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.GET("/api/farms", handler.GetFarms)
	router.NoRoute(handler.ResponseBadRequest)

	//пропишем пользователя
	router.POST("/api/user/register", handler.Register)
	for _, tt := range reliableData.users {
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}
	//пропишем фермы
	router.POST("/api/farms", handler.AddFarm)
	for _, tt := range reliableData.farms {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{Name: config.Cookie, Value: tt.cookie, HttpOnly: true}
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/farms", body)
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
	}

	//запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.args.cookie, HttpOnly: true}
			req, err := http.NewRequest("GET", "/api/farms", nil)
			req.AddCookie(cookie)

			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			input, _ := ioutil.ReadAll(w.Body)
			assert.Equal(t, tt.want.body, string(input))
		})
	}
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

func TestBaseHandler_GetCowBreeds(t *testing.T) {
	type user struct {
		body string
	}
	reliableData := struct {
		user user
	}{
		user: user{
			body: "{\"login\":\"User\",\"password\":\"pa$$word_1\"}",
		},
	}

	type want struct {
		code int
		body string
	}
	type args struct {
		cookie string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get breeds",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			want: want{
				code: http.StatusOK,
				body: "[{\"breed_id\":1,\"breed\":\"Голштинская\"},{\"breed_id\":2,\"breed\":\"Красная датская\"},{\"breed_id\":3,\"breed\":\"Айрширская\"}]",
			},
		},
		{
			name: "No user",
			args: args{
				cookie: "a09bccf2b294353452b34dc0e08d8b582a",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "{\"Message\":\"Unauthorized\"}",
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.GET("/api/cows/breeds", handler.GetCowBreeds)
	router.NoRoute(handler.ResponseBadRequest)

	//пропишем пользователя
	router.POST("/api/user/register", handler.Register)
	w := httptest.NewRecorder()
	body := bytes.NewBuffer([]byte(reliableData.user.body))
	req, _ := http.NewRequest("POST", "/api/user/register", body)
	router.ServeHTTP(w, req)

	//запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.args.cookie, HttpOnly: true}
			req, err := http.NewRequest("GET", "/api/cows/breeds", nil)
			req.AddCookie(cookie)

			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			input, _ := ioutil.ReadAll(w.Body)
			assert.Equal(t, tt.want.body, string(input))
		})
	}
}

func TestBaseHandler_AddCow(t *testing.T) {
	type user struct {
		body string
	}
	type farm struct {
		body   string
		cookie string
	}
	reliableData := struct {
		users []user
		farms []farm
	}{
		users: []user{
			{
				body: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			},
			{
				body: "{\"login\": \"User2\", \"password\": \"pa$$word_2\"}",
			},
		},
		farms: []farm{
			{
				body:   "{\"name\": \"Farm1\",\"address\": \"Moscow\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"Farm2\",\"address\": \"Omsk\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
	}

	type want struct {
		code int
		body string
	}
	type args struct {
		cookie string
		body   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Add cow",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "{\"name\": \"корова 1\",\"breed_id\": 1,\"farm_id\": 1,\"bolus_sn\": 1221,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
			},
			want: want{
				code: http.StatusCreated,
				body: "",
			},
		},
		{
			name: "Duplicate bolus",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "{\"name\": \"корова 1\",\"breed_id\": 1,\"farm_id\": 2,\"bolus_sn\": 1221,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
			},
			want: want{
				code: http.StatusConflict,
				body: "",
			},
		},
		{
			name: "No user",
			args: args{
				cookie: "a09bccf2b294353452b34dc0e08d8b582a",
				body:   "{\"name\": \"корова 2\",\"breed_id\": 1,\"farm_id\": 1,\"bolus_sn\": 1223,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "{\"Message\":\"Unauthorized\"}",
			},
		},
		{
			name: "Bad request",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "wrong json",
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.POST("/api/cows", handler.AddCow)
	router.NoRoute(handler.ResponseBadRequest)

	//пропишем пользователя
	router.POST("/api/user/register", handler.Register)
	for _, tt := range reliableData.users {
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}
	//пропишем фермы
	router.POST("/api/farms", handler.AddFarm)
	for _, tt := range reliableData.farms {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{Name: config.Cookie, Value: tt.cookie, HttpOnly: true}
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/farms", body)
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.args.cookie, HttpOnly: true}

			body := bytes.NewBuffer([]byte(tt.args.body))
			req, err := http.NewRequest("POST", "/api/cows", body)
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

func TestBaseHandler_DelCows(t *testing.T) {
	type user struct {
		body string
	}
	type farm struct {
		body   string
		cookie string
	}
	type cow struct {
		body   string
		cookie string
	}
	reliableData := struct {
		users []user
		farms []farm
		cows  []cow
	}{
		users: []user{
			{
				body: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			},
			{
				body: "{\"login\": \"User2\", \"password\": \"pa$$word_2\"}",
			},
		},
		farms: []farm{
			{
				body:   "{\"name\": \"Farm1\",\"address\": \"Moscow\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"Farm2\",\"address\": \"Omsk\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
		cows: []cow{
			{
				body:   "{\"name\": \"корова 1\",\"breed_id\": 1,\"farm_id\": 1,\"bolus_sn\": 1221,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"корова 2\",\"breed_id\": 1,\"farm_id\": 1,\"bolus_sn\": 1222,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"корова 3\",\"breed_id\": 1,\"farm_id\": 1,\"bolus_sn\": 1223,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
	}

	type want struct {
		code int
		body string
	}
	type args struct {
		cookie string
		body   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Del cow",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "[1,2,3]",
			},
			want: want{
				code: http.StatusAccepted,
				body: "",
			},
		},
		{
			name: "Missing cows",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "[1]",
			},
			want: want{
				code: http.StatusConflict,
				body: "",
			},
		},
		{
			name: "No user",
			args: args{
				cookie: "a09bccf2b294353452b34dc0e08d8b582a",
				body:   "[1]",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "{\"Message\":\"Unauthorized\"}",
			},
		},
		{
			name: "Bad request",
			args: args{
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
				body:   "wrong json",
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}

	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.DELETE("/api/cows", handler.DelCows)
	router.NoRoute(handler.ResponseBadRequest)

	//пропишем пользователя
	router.POST("/api/user/register", handler.Register)
	for _, tt := range reliableData.users {
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}
	//пропишем фермы
	router.POST("/api/farms", handler.AddFarm)
	for _, tt := range reliableData.farms {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{Name: config.Cookie, Value: tt.cookie, HttpOnly: true}
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/farms", body)
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
	}
	//пропишем коров
	router.POST("/api/cows", handler.AddCow)
	for _, tt := range reliableData.cows {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{Name: config.Cookie, Value: tt.cookie, HttpOnly: true}
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/cows", body)
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.args.cookie, HttpOnly: true}

			body := bytes.NewBuffer([]byte(tt.args.body))
			req, err := http.NewRequest("DELETE", "/api/cows", body)
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

func TestBaseHandler_GetCows(t *testing.T) {
	type user struct {
		body string
	}
	type farm struct {
		body   string
		cookie string
	}
	type cow struct {
		body   string
		cookie string
	}

	reliableData := struct {
		users []user
		farms []farm
		cows  []cow
	}{
		users: []user{
			{
				body: "{\"login\": \"User\", \"password\": \"pa$$word_1\"}",
			},
			{
				body: "{\"login\": \"User2\", \"password\": \"pa$$word_2\"}",
			},
		},
		farms: []farm{
			{
				body:   "{\"name\": \"Farm1\",\"address\": \"Moscow\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"Farm2\",\"address\": \"Omsk\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
		cows: []cow{
			{
				body:   "{\"name\": \"корова 1\",\"breed_id\": 1,\"farm_id\": 1,\"bolus_sn\": 1221,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"корова 2\",\"breed_id\": 1,\"farm_id\": 1,\"bolus_sn\": 1222,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			{
				body:   "{\"name\": \"корова 3\",\"breed_id\": 1,\"farm_id\": 2,\"bolus_sn\": 1223,\"date_of_born\": \"2020-04-09T23:00:00Z\",\"bolus_type\": \"С датчиком PH\"}",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
		},
	}
	type data struct {
		bolusNum int
		farmID   int
		name     string
	}
	type want struct {
		code int
		body string
		data []data
	}
	type args struct {
		farm   string
		cookie string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get cow farm1",
			args: args{
				farm:   "1",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			want: want{
				code: http.StatusOK,
				data: []data{
					{
						bolusNum: 1221,
						farmID:   1,
						name:     "корова 1",
					},
					{
						bolusNum: 1222,
						farmID:   1,
						name:     "корова 2",
					},
				},
			},
		},
		{
			name: "No farm",
			args: args{
				farm:   "10",
				cookie: "8f9bfe9d1345237cb3b2b205864da075",
			},
			want: want{
				code: http.StatusNoContent,
				body: "",
			},
		},
		{
			name: "No user",
			args: args{
				farm:   "1",
				cookie: "8f9bfe3249d1345237cb3b2b205864da075",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "{\"Message\":\"Unauthorized\"}",
			},
		},
	}
	logger.Wr = logger.New()
	ds := storage.NewDBMockStorage(context.Background())
	handler := NewBaseHandler(ds)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.GET("/api/farms/:id/cows", handler.GetCows)
	router.NoRoute(handler.ResponseBadRequest)

	//пропишем пользователя
	router.POST("/api/user/register", handler.Register)
	for _, tt := range reliableData.users {
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/user/register", body)
		router.ServeHTTP(w, req)
	}
	//пропишем фермы
	router.POST("/api/farms", handler.AddFarm)
	for _, tt := range reliableData.farms {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{Name: config.Cookie, Value: tt.cookie, HttpOnly: true}
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/farms", body)
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
	}
	//пропишем коров
	router.POST("/api/cows", handler.AddCow)
	for _, tt := range reliableData.cows {
		w := httptest.NewRecorder()
		cookie := &http.Cookie{Name: config.Cookie, Value: tt.cookie, HttpOnly: true}
		body := bytes.NewBuffer([]byte(tt.body))
		req, _ := http.NewRequest("POST", "/api/cows", body)
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
	}

	//запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			cookie := &http.Cookie{Name: config.Cookie, Value: tt.args.cookie, HttpOnly: true}
			req, err := http.NewRequest("GET", "/api/farms/"+tt.args.farm+"/cows", nil)
			req.AddCookie(cookie)

			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			input, _ := ioutil.ReadAll(w.Body)
			if w.Code == http.StatusOK {
				var cowResp []storage.Cow
				err = json.Unmarshal(input, &cowResp)
				assert.Nil(t, err)
				for i, c := range cowResp {
					assert.Equal(t, c.BolusNum, tt.want.data[i].bolusNum)
					assert.Equal(t, c.FarmID, tt.want.data[i].farmID)
					assert.Equal(t, c.Name, tt.want.data[i].name)
				}
			} else {
				assert.Equal(t, tt.want.body, string(input))
			}
		})
	}
}

func Test_getIDFromJSON(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []int
		wantErr bool
	}{
		{
			name: "one entry",
			args: "[1]",
			want: []int{
				1,
			},
			wantErr: false,
		},
		{
			name: "many entries",
			args: "[3,5,8]",
			want: []int{
				3,
				5,
				8,
			},
			wantErr: false,
		},
		{
			name:    "wrong",
			args:    "[w1]",
			wantErr: true,
		},
	}
	logger.Wr = logger.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := ioutil.NopCloser(strings.NewReader(tt.args))
			got, err := getIDFromJSON(in)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, got, tt.want)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
