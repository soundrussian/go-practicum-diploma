package v1

import (
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		secretKey *string
		storage   storage.Storage
	}
	tests := []struct {
		name    string
		args    args
		want    *Auth
		wantErr bool
	}{
		{
			name: "returns error if secretKey is nil",
			args: args{
				storage: new(mocks.Storage),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns error if storage is nil",
			args: args{
				secretKey: secretKey,
				storage:   nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns initialized auth service with passed storage",
			args: args{
				secretKey: secretKey,
				storage:   new(mocks.Storage),
			},
			want:    &Auth{storage: new(mocks.Storage)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretKey = tt.args.secretKey
			got, err := New(tt.args.storage)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
