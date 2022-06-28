package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoundTime(t *testing.T) {
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1 year 1 month",
			args: args{
				from: time.Date(2015, 5, 0, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2016, 6, 0, 0, 0, 0, 0, time.UTC),
			},
			want: "1 г 1 мес",
		},
		{
			name: "0 year 5 month",
			args: args{
				from: time.Date(2020, 5, 0, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2020, 10, 1, 0, 30, 0, 0, time.UTC),
			},
			want: "5 мес",
		},
		{
			name: "10 year 0 month",
			args: args{
				from: time.Date(2010, 10, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2020, 10, 1, 0, 30, 0, 0, time.UTC),
			},
			want: "10 л",
		},
		{
			name: "23 year 11 month",
			args: args{
				from: time.Date(2000, 0, 5, 0, 0, 10, 0, time.UTC),
				to:   time.Date(2023, 11, 5, 0, 30, 0, 0, time.UTC),
			},
			want: "23 г 11 мес",
		},
		{
			name: "reverse 15 year 3 month",
			args: args{
				from: time.Date(2010, 7, 5, 8, 34, 10, 0, time.UTC),
				to:   time.Date(1995, 4, 9, 0, 30, 0, 0, time.UTC),
			},
			want: "15 л 2 мес",
		},
		{
			name: "3 year 6 month",
			args: args{
				from: time.Date(2000, 11, 5, 8, 34, 10, 0, time.UTC),
				to:   time.Date(2004, 5, 9, 0, 30, 0, 0, time.UTC),
			},
			want: "3 г 6 мес",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundTime(tt.args.from, tt.args.to)
			assert.Equal(t, tt.want, got)
		})
	}
}
