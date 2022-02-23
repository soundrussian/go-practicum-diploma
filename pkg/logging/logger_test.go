package logging

import (
	"github.com/rs/zerolog"
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		opts []LoggerOption
	}
	tests := []struct {
		name string
		args args
		want zerolog.Level
	}{
		{
			name: "returns new logger with default (trace) level if no options are passed",
			args: args{opts: nil},
			want: zerolog.TraceLevel,
		},
		{
			name: "returns new logger with level set by passed level option",
			args: args{opts: []LoggerOption{WithLogLevel(zerolog.ErrorLevel)}},
			want: zerolog.ErrorLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(tt.args.opts...).GetLevel(); got != tt.want {
				t.Errorf("NewLogger().Level() = %v, want %v", got, tt.want)
			}
		})
	}
}
