package api

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPIConfig_RunAddress(t *testing.T) {
	type args struct {
		Env  string
		Flag string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sets defaultRunAddress by default",
			args: args{},
			want: defaultRunAddress,
		},
		{
			name: "sets value from env if it is set",
			args: args{
				Env: "localhost:3333",
			},
			want: "localhost:3333",
		},
		{
			name: "sets value from flag if it is set",
			args: args{
				Flag: "localhost:4444",
			},
			want: "localhost:4444",
		},
		{
			name: "sets value from flag if both flag and env are set",
			args: args{
				Env:  "localhost:3333",
				Flag: "localhost:4444",
			},
			want: "localhost:4444",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.Env != "" {
				t.Setenv(runAddressEnvKey, tt.args.Env)
			}
			if tt.args.Flag != "" {
				err := flag.Set(runAddressFlagName, tt.args.Flag)
				require.NoError(t, err)
			}
			flag.Parse()
			readConfig()

			assert.Equal(t, tt.want, config.RunAddress)
		})
	}
}
