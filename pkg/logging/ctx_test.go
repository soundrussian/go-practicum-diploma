package logging

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCtxLogger_StoresLogger(t *testing.T) {
	t.Run("stores logger in the passed context", func(t *testing.T) {
		ctx, logger := CtxLogger(nil)

		// Fetch logger again from the context
		_, got := CtxLogger(ctx)

		assert.Equal(t, logger, got)
	})
}

func TestCtxLogger_CorrelationID(t *testing.T) {
	ctxWithLogger, _ := CtxLogger(nil)

	type args struct {
		ctx           context.Context
		correlationID string
	}
	tests := []struct {
		name        string
		args        args
		wantExactly string
		wantErr     bool
	}{
		{
			name:    "stores new correlation id if passed context is nil",
			args:    args{},
			wantErr: false,
		},
		{
			name: "keeps correlation id if it is already set in context",
			args: args{
				ctx:           ctxWithLogger,
				correlationID: "existing-correlation-id",
			},
			wantErr:     false,
			wantExactly: "existing-correlation-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.args.ctx
			if tt.args.correlationID != "" {
				ctx = context.WithValue(ctx, contextKeyCorrelationID, tt.args.correlationID)
			}

			ctx, _ = CtxLogger(ctx)
			got, err := CorrelationID(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.wantExactly != "" {
				assert.Equal(t, tt.wantExactly, got)
			}
		})
	}
}

func TestSetCorrelationID(t *testing.T) {
	ctxWithLogger, _ := CtxLogger(nil)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "overwrites corellation id in context",
			args: args{
				ctx: ctxWithLogger,
				id:  "my-correlation-id",
			},
			want: "my-correlation-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := SetCorrelationID(tt.args.ctx, tt.args.id)
			got, err := CorrelationID(ctx)
			require.NoError(t, err)
			assert.Equal(t, tt.args.id, got)
		})
	}
}

func TestCorrelationID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "returns correlation id if it is set in context",
			args: args{
				ctx: context.WithValue(context.Background(), contextKeyCorrelationID, "test-correlation-id"),
			},
			want:    "test-correlation-id",
			wantErr: assert.NoError,
		},
		{
			name: "returns error if there is no correlation id in the context",
			args: args{
				ctx: context.Background(),
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CorrelationID(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("CorrelationID(%v)", tt.args.ctx)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CorrelationID(%v)", tt.args.ctx)
		})
	}
}
