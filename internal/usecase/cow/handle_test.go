package cow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/Sur0vy/cows_health.git/errors"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/logger"
	storageMock "github.com/Sur0vy/cows_health.git/mocks"
)

func TestHandler_Add(t *testing.T) {
	const dateFormat = "2006-01-02"
	type args struct {
		useMocks bool
		cow      models.Cow
		dt       string
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
			name: "cow already exist",
			args: args{
				useMocks: true,
				cow: models.Cow{
					Name:     "Cow_1",
					BreedID:  1,
					FarmID:   1,
					BolusNum: 1,
				},
				dt:  "2022-02-01",
				err: errors.ErrExist,
			},
			want: want{
				code: http.StatusConflict,
			},
		},
		{
			name: "cow already exist",
			args: args{
				useMocks: true,
				cow: models.Cow{
					Name:     "Cow_1",
					BreedID:  1,
					FarmID:   1,
					BolusNum: 1,
				},
				dt:  "2022-02-01",
				err: errors.ErrEmpty,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "bad body",
			args: args{
				useMocks: false,
				cow: models.Cow{
					Name:     "Cow_1",
					BreedID:  1,
					FarmID:   1,
					BolusNum: 1,
				},
				badBody: "bad body",
				dt:      "2022-02-01",
				err:     nil,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	repo := &storageMock.CowStorage{}

	ch := NewCowHandler(repo, logger.New())
	router := echo.New()
	router.POST("/api/cow", ch.Add)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.args.useMocks {
				tt.args.cow.DateOfBorn, _ = time.Parse(dateFormat, tt.args.dt)
				tt.args.cow.AddedAt, _ = time.Parse(dateFormat, tt.args.dt)
				repo.On("Add", context.Background(), tt.args.cow).
					Return(tt.args.err).
					Once()
				body, _ = json.Marshal(tt.args.cow)
			} else {
				body = []byte(tt.args.badBody)
			}
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/cow", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type args struct {
		useMocks bool
		body     string
		arr      []int
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
			name: "success",
			args: args{
				useMocks: true,
				body:     "[1, 2, 3]",
				arr:      []int{1, 2, 3},
				err:      nil,
			},
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: "bad body",
			args: args{
				useMocks: false,
				body:     "[1, 2, ",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "empty",
			args: args{
				useMocks: true,
				body:     "[4]",
				arr:      []int{4},
				err:      errors.ErrEmpty,
			},
			want: want{
				code: http.StatusConflict,
			},
		},
		{
			name: "unexpected",
			args: args{
				useMocks: true,
				body:     "[1, 2, 3]",
				arr:      []int{1, 2, 3},
				err:      errors.ErrExist,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	repo := &storageMock.CowStorage{}

	ch := NewCowHandler(repo, logger.New())
	router := echo.New()
	router.DELETE("/api/cow", ch.Delete)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.useMocks {
				repo.On("Delete", context.Background(), tt.args.arr).
					Return(tt.args.err).
					Once()
			}
			recorder := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(tt.args.body))
			req, err := http.NewRequest("DELETE", "/api/cow", body)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, req)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.code, recorder.Code)
		})
	}
}

func TestHandler_Get(t *testing.T) {
	const dateFormat = "2006-01-02"
	type args struct {
		useMocks bool
		id       string
		cows     []models.Cow
		dt       string
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
			name: "success",
			args: args{
				useMocks: true,
				id:       "1",
				cows: []models.Cow{
					{
						Name:     "Cow_1",
						ID:       1,
						BreedID:  1,
						FarmID:   1,
						BolusNum: 1,
					},
				},
				dt:  "2022-02-01",
				err: nil,
			},
			want: want{
				code: http.StatusOK,
				body: "[{\"id\":1,\"name\":\"Cow_1\",\"breed_id\":1,\"farm_id\":1,\"bolus_sn\":1,\"date_of_born\":\"2022-02-01T00:00:00Z\",\"added_at\":\"2022-02-01T00:00:00Z\"}]\n",
			},
		},
		{
			name: "empty",
			args: args{
				useMocks: true,
				id:       "1",
				err:      errors.ErrEmpty,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "unexpected",
			args: args{
				useMocks: true,
				id:       "1",
				err:      errors.ErrExist,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "success many",
			args: args{
				useMocks: true,
				id:       "1",
				cows: []models.Cow{
					{
						Name:     "Cow_1",
						ID:       1,
						BreedID:  1,
						FarmID:   1,
						BolusNum: 1,
					},
					{
						Name:     "Cow_2",
						ID:       2,
						BreedID:  2,
						FarmID:   1,
						BolusNum: 2,
					},
					{
						Name:     "Cow_3",
						ID:       3,
						BreedID:  1,
						FarmID:   1,
						BolusNum: 3,
					},
				},
				dt:  "2022-02-01",
				err: nil,
			},
			want: want{
				code: http.StatusOK,
				body: "[{\"id\":1,\"name\":\"Cow_1\",\"breed_id\":1,\"farm_id\":1,\"bolus_sn\":1,\"date_of_born\":\"2022-02-01T00:00:00Z\",\"added_at\":\"2022-02-01T00:00:00Z\"}," +
					"{\"id\":2,\"name\":\"Cow_2\",\"breed_id\":2,\"farm_id\":1,\"bolus_sn\":2,\"date_of_born\":\"2022-02-01T00:00:00Z\",\"added_at\":\"2022-02-01T00:00:00Z\"}," +
					"{\"id\":3,\"name\":\"Cow_3\",\"breed_id\":1,\"farm_id\":1,\"bolus_sn\":3,\"date_of_born\":\"2022-02-01T00:00:00Z\",\"added_at\":\"2022-02-01T00:00:00Z\"}]\n",
			},
		},
		{
			name: "bad",
			args: args{
				useMocks: false,
				id:       "bad",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	repo := &storageMock.CowStorage{}

	ch := NewCowHandler(repo, logger.New())
	router := echo.New()
	router.GET("/api/cow/:id", ch.Get)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.useMocks {
				for i := range tt.args.cows {
					tt.args.cows[i].DateOfBorn, _ = time.Parse(dateFormat, tt.args.dt)
					tt.args.cows[i].AddedAt, _ = time.Parse(dateFormat, tt.args.dt)
				}
				id, _ := strconv.Atoi(tt.args.id)
				repo.On("Get", context.Background(), id).
					Return(tt.args.cows, tt.args.err).
					Once()
			}
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/cow/%s", tt.args.id)
			req, err := http.NewRequest("GET", url, nil)
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

func TestHandler_GetBreeds(t *testing.T) {
	type args struct {
		breeds []models.Breed
		err    error
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
			name: "success",
			args: args{
				breeds: []models.Breed{
					{
						1,
						"Порода 1",
					},
					{
						2,
						"Порода 2",
					},
				},
				err: nil,
			},
			want: want{
				code: http.StatusOK,
				body: "[{\"breed_id\":1,\"breed\":\"Порода 1\"},{\"breed_id\":2,\"breed\":\"Порода 2\"}]\n",
			},
		},
		{
			name: "empty",
			args: args{
				err: errors.ErrEmpty,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "bad database request",
			args: args{
				err: errors.ErrExist,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	repo := &storageMock.CowStorage{}

	ch := NewCowHandler(repo, logger.New())
	router := echo.New()
	router.GET("/api/cow", ch.GetBreeds)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetBreeds", context.Background()).
				Return(tt.args.breeds, tt.args.err).
				Once()

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/cow", nil)
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

func TestHandler_GetInfo(t *testing.T) {
	const dateFormat = "2006-01-02"
	type args struct {
		useMocks bool
		id       string
		info     models.CowInfo
		dt       string
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
			name: "success",
			args: args{
				useMocks: true,
				id:       "1",
				info: models.CowInfo{
					Health: models.Health{
						Ill:    "Здорова",
						Estrus: false,
					},
					Summary: models.Cow{
						Name:     "Корова_1",
						Breed:    "вологодская",
						BolusNum: 1,
					},
					History: []models.MonitoringData{
						{
							PH:          6.7,
							Temperature: 53,
							Charge:      30,
							Movement:    24,
						},
						{
							PH:          6.5,
							Temperature: 50,
							Charge:      60,
							Movement:    15,
						},
					},
				},
				dt:  "2022-02-01",
				err: nil,
			},
			want: want{
				code: http.StatusOK,
				body: "{\"health\":" +
					"{\"ill\":\"Здорова\",\"estrus\":false,\"updated_at\":\"2022-02-01T00:00:00Z\"}," +
					"\"summary\":" +
					"{\"name\":\"Корова_1\",\"breed\":\"вологодская\",\"bolus_sn\":1,\"date_of_born\":\"2022-02-01T00:00:00Z\",\"added_at\":\"2022-02-01T00:00:00Z\"}," +
					"\"history\":[" +
					"{\"added_at\":\"2022-02-01T00:00:00Z\",\"ph\":6.7,\"temperature\":53,\"movement\":24,\"charge\":30}," +
					"{\"added_at\":\"2022-02-01T00:00:00Z\",\"ph\":6.5,\"temperature\":50,\"movement\":15,\"charge\":60}]}\n",
			},
		},
		{
			name: "bad",
			args: args{
				useMocks: false,
				id:       "bad",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "empty",
			args: args{
				useMocks: true,
				id:       "1",
				err:      errors.ErrEmpty,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "empty",
			args: args{
				useMocks: true,
				id:       "1",
				err:      errors.ErrExist,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	repo := &storageMock.CowStorage{}

	ch := NewCowHandler(repo, logger.New())
	router := echo.New()
	router.GET("/api/cow/info/:id", ch.GetInfo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.useMocks {
				tt.args.info.Health.UpdatedAt, _ = time.Parse(dateFormat, tt.args.dt)
				tt.args.info.Summary.DateOfBorn, _ = time.Parse(dateFormat, tt.args.dt)
				tt.args.info.Summary.AddedAt, _ = time.Parse(dateFormat, tt.args.dt)
				for i := range tt.args.info.History {
					tt.args.info.History[i].AddedAt, _ = time.Parse(dateFormat, tt.args.dt)
				}
				id, _ := strconv.Atoi(tt.args.id)
				repo.On("GetInfo", context.Background(), id).
					Return(tt.args.info, tt.args.err).
					Once()
			}
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/cow/info/%s", tt.args.id)
			req, err := http.NewRequest("GET", url, nil)
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

func Test_getIDFromJSON(t *testing.T) {
	type args struct {
		body string
	}
	type want struct {
		arr []int
		err bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success",
			args: args{
				body: "[1, 2, 3]",
			},
			want: want{
				arr: []int{1, 2, 3},
				err: false,
			},
		},
		{
			name: "success 2",
			args: args{
				body: "[399]",
			},
			want: want{
				arr: []int{399},
				err: false,
			},
		},
		{
			name: "fail",
			args: args{
				body: "bad input",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ioutil.NopCloser(strings.NewReader(tt.args.body)) // r type is io.ReadCloser
			out, err := getIDFromJSON(r)
			if tt.want.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want.arr, out)
			}
		})
	}
}
