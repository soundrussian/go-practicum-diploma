package v1

import (
	"github.com/soundrussian/go-practicum-diploma/storage"
	"github.com/soundrussian/go-practicum-diploma/storage/mock"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		name    string
		args    args
		want    *Order
		wantErr bool
	}{
		{
			name:    "returns error if passed storage is nil",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns Order with passed storage if it is not nil",
			args: args{
				storage: new(mock.Storage),
			},
			want:    &Order{storage: new(mock.Storage)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
