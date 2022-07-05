package farm

import (
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
)

func TestHandler_Add(t *testing.T) {
	type fields struct {
		log         *logger.Logger
		farmStorage models.FarmStorage
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
				log:         tt.fields.log,
				farmStorage: tt.fields.farmStorage,
			}
			if err := h.Add(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type fields struct {
		log         *logger.Logger
		farmStorage models.FarmStorage
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
				log:         tt.fields.log,
				farmStorage: tt.fields.farmStorage,
			}
			if err := h.Delete(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	type fields struct {
		log         *logger.Logger
		farmStorage models.FarmStorage
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
				log:         tt.fields.log,
				farmStorage: tt.fields.farmStorage,
			}
			if err := h.Get(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewFarmHandler(t *testing.T) {
	type args struct {
		fs  models.FarmStorage
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
			if got := NewFarmHandler(tt.args.fs, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFarmHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
