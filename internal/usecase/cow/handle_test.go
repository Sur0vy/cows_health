package cow

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	storageMock "github.com/Sur0vy/cows_health.git/internal/mocks"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestHandler_Add(t *testing.T) {
	const dateFormat = "2006-01-02"
	type args struct {
		useMocks bool
		cow      models.Cow
		dt       string
		//body     string
		err error
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
			name: "success",
			args: args{
				useMocks: true,
				cow: models.Cow{
					Name:     "Cow_1",
					BreedID:  1,
					FarmID:   1,
					BolusNum: 1,
				},
				dt:  "2022-02-01",
				err: nil,
			},
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "bad body",
		},
		{
			name: "exist",
		},
		{
			name: "unexpected error",
		},
	}

	repo := &storageMock.CowStorage{}

	ch := NewCowHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/cow", ch.Add)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.useMocks {
				tt.args.cow.DateOfBorn, _ = time.Parse(dateFormat, tt.args.dt)
				tt.args.cow.AddedAt, _ = time.Parse(dateFormat, tt.args.dt)
				repo.On("Add", context.Background(), tt.args.cow).
					Return(tt.args.err).
					Once()
			}
			body, _ := json.Marshal(tt.args.cow)
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/cow", bytes.NewReader(body))
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
	type fields struct {
		log *logger.Logger
		cs  models.CowStorage
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
				cs:  tt.fields.cs,
			}
			if err := h.Delete(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	type fields struct {
		log *logger.Logger
		cs  models.CowStorage
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
				cs:  tt.fields.cs,
			}
			if err := h.Get(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_GetBreeds(t *testing.T) {
	type fields struct {
		log *logger.Logger
		cs  models.CowStorage
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
				cs:  tt.fields.cs,
			}
			if err := h.GetBreeds(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetBreeds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_GetInfo(t *testing.T) {
	type fields struct {
		log *logger.Logger
		cs  models.CowStorage
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
				cs:  tt.fields.cs,
			}
			if err := h.GetInfo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCowHandler(t *testing.T) {
	type args struct {
		cs  models.CowStorage
		log *logger.Logger
	}
	tests := []struct {
		name string
		args args
		want Handle
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCowHandler(tt.args.cs, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCowHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIDFromJSON(t *testing.T) {
	type args struct {
		reader io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIDFromJSON(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIDFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIDFromJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
