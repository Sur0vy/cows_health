package handlers

import (
	"bytes"
	"context"
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
	//type fields struct {
	//	body []string
	//}
	//type args struct {
	//	id string
	//}
	//type want struct {
	//	URL  string
	//	code int
	//}
	//tests := []struct {
	//	name   string
	//	fields fields
	//	args   args
	//	want   want
	//}{
	//	{
	//		name: "Test GET correct",
	//		fields: fields{
	//			body: []string{"testURL1"},
	//		},
	//		args: args{
	//			id: "testShortURL1",
	//		},
	//		want: want{
	//			URL:  "testURL1",
	//			code: http.StatusTemporaryRedirect,
	//		},
	//	},
	//}
	//
	//for _, tt := range tests {
	//t.Run(tt.name, func(t *testing.T) {
	//		storage := storage.NewDBStorage(context.Background(), "")
	//		s := server.SetupServer(storage)
	//		w := httptest.NewRecorder()
	//
	//		//заполним БД
	//		for _, v := range tt.fields.body {
	//			body := bytes.NewBuffer([]byte(v))
	//			req, _ := http.NewRequest("POST", "/", body)
	//			s.ServeHTTP(w, req)
	//		}
	//
	//		r := httptest.NewRecorder()
	//
	//		URL := "/" + tt.args.id
	//		req, err := http.NewRequest("GET", URL, nil)
	//		s.ServeHTTP(r, req)
	//
	//		assert.Nil(t, err)
	//		assert.Equal(t, tt.want.code, r.Code)
	//
	//		URL = r.Header().Get("Location")
	//		assert.Equal(t, tt.want.URL, URL)
	//		}
	//}
}

func TestBaseHandler_DelFarm(t *testing.T) {

}

func TestBaseHandler_GetFarms(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//
	//s := mocks.NewMockStorage(ctrl)
	//
	//value := []byte("Some value")
	//
	//// при вызове с произвольным аргументом
	//// заглушка будет возвращать слайс
	//// метод может быть вызван не более 5 раз
	//s.EXPECT().
	//	GetFarms("1").
	//	Return(value, nil).
	//	MaxTimes(5)
	//
	//// тестируем функцию
	//Lookup(s, someCond)
}

func TestBaseHandler_Login(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//
	//storage := mocks.NewMockStorage(ctrl)
	//
	//s := server.SetupServer(storage)
	//w := httptest.NewRecorder()

	//body := bytes.NewBuffer([]byte(tt.args.body))
	//req, err := http.NewRequest("POST", "/", body)
	//
	//s.ServeHTTP(w, req)
	//
	//assert.Nil(t, err)
	//
	//assert.Equal(t, w.Code, tt.want.code)
	//
	//if tt.args.trueVal {
	//	body = bytes.NewBuffer([]byte(config.HTTP + config.Cnf.BaseURL + "/" + tt.want.body))
	//	assert.Equal(t, fmt.Sprint(body), fmt.Sprint(w.Body))
	//} else {
	//	body = bytes.NewBuffer([]byte(config.Cnf.BaseURL + tt.want.body))
	//	assert.NotEqual(t, fmt.Sprint(w.Body), fmt.Sprint(body))
	//}

	//value := []byte("Some value")

	// при вызове с произвольным аргументом
	// заглушка будет возвращать слайс
	// метод может быть вызван не более 5 раз
	//storage.EXPECT().
	//	Login("1").
	//	Return(value, nil).
	//	MaxTimes(5)

	// тестируем функцию
	//Lookup(s, someCond)
}

func TestBaseHandler_Logout(t *testing.T) {

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
	api := router.Group("/api")
	user := api.Group("/user")
	user.POST("/register", handler.Register)
	router.NoRoute(handler.ResponseBadRequest)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(tt.args))
			req, err := http.NewRequest("POST", "/api/user/register", body)
			router.ServeHTTP(w, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestBaseHandler_ResponseBadRequest(t *testing.T) {

}

func TestNewBaseHandler(t *testing.T) {

}

func Test_getIDFromJSON(t *testing.T) {

}
