package helpers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_checkPassword(t *testing.T) {
	type args struct {
		hash     string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "right",
			args: args{
				hash:     "24326124313024616d4c5137644b4839396f6167685361484f31496765316461625438706d512e74493656557142544643482f2e37765a4d2f6a5175",
				password: "pa$$word_1",
			},
			want: true,
		},
		{
			name: "wrong",
			args: args{
				hash:     "24326124313024616d4c5137644b22139396f6167685361484f31496765316461625438706d512e74493656557142544643482f2e37765a4d2f6a5175",
				password: "pa$$word",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.args.hash, tt.args.password); got != tt.want {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func Test_getCryptoPassword(t *testing.T) {
	type args struct {
		password   string
		passwordCh string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wrong   bool
	}{
		{
			name: "right",
			args: args{
				password:   "pa$$word_1",
				passwordCh: "pa$$word_1",
			},
			wantErr: false,
			wrong:   false,
		},
		{
			name: "wrong",
			args: args{
				password:   "pa$$word",
				passwordCh: "pa$$word_1",
			},
			wantErr: false,
			wrong:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GetCryptoPassword(tt.args.password)
			if !tt.wantErr {
				res := CheckPassword(hash, tt.args.passwordCh)
				assert.NotEqual(t, res, tt.wrong)
				require.NoError(t, err)
				return
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func Test_getMD5Hash(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			name:    "right",
			args:    "user",
			want:    "ee11cbb19052e40b07aac0ca060c23ee",
			wantErr: false,
		},
		{
			name:    "wrong",
			args:    "user",
			want:    "8f9bfe9d1345237cb3b2b205864da075",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMD5Hash(tt.args)
			if !tt.wantErr {
				assert.Equal(t, got, tt.want)
				return
			}
			assert.NotEqual(t, got, tt.want)
		})
	}
}
