package v1

import (
	"context"
	"errors"
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"github.com/stretchr/testify/mock"
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
		want    *Balance
		wantErr bool
	}{
		{
			name:    "returns error if passed storage is nil",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns initialized service with storage set",
			args: args{
				storage: new(mocks.Storage),
			},
			want:    &Balance{storage: new(mocks.Storage)},
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

func TestBalance_UserBalance(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserBalance
		wantErr bool
	}{
		{
			name: "returns nil and err if storage failed to get user balance",
			fields: fields{
				storage: failingStorage(),
			},
			args: args{
				userID: 100,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns balance from storage and no error if storage does not return error",
			fields: fields{
				storage: successfulStorage(100, 50),
			},
			args: args{
				userID: 100,
			},
			want:    &model.UserBalance{Current: 100, Withdrawn: 50},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balance{
				storage: tt.fields.storage,
			}
			got, err := b.UserBalance(context.Background(), tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func failingStorage() *mocks.Storage {
	m := new(mocks.Storage)
	m.On("UserBalance", mock.Anything, mock.Anything).Return(nil, errors.New("mock error"))
	return m
}

func successfulStorage(current uint64, withdrawn uint64) *mocks.Storage {
	m := new(mocks.Storage)
	m.On("UserBalance", mock.Anything, mock.Anything).Return(
		&model.UserBalance{
			Current:   current,
			Withdrawn: withdrawn,
		},
		nil,
	)
	return m
}

func TestBalance_Withdraw(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		userID     uint64
		withdrawal model.Withdrawal
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "returns error if sum is less than zero",
			fields: fields{
				storage: new(mocks.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Sum: -1},
			},
			wantErr: true,
		},
		{
			name: "returns error if sum is zero",
			fields: fields{
				storage: new(mocks.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Sum: 0},
			},
			wantErr: true,
		},
		{
			name: "returns error if order number is missing",
			fields: fields{
				storage: new(mocks.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Sum: 10},
			},
			wantErr: true,
		},
		{
			name: "returns error if order number is not a number",
			fields: fields{
				storage: new(mocks.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "not a number", Sum: 10},
			},
			wantErr: true,
		},
		{
			name: "returns error if order checksum is invalid",
			fields: fields{
				storage: new(mocks.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "7992739871", Sum: 10},
			},
			wantErr: true,
		},
		{
			name: "returns error if storage reported error",
			fields: fields{
				storage: failingWithdrawal(),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "79927398713", Sum: 10},
			},
			wantErr: true,
		},
		{
			name: "does not return error if storage reported success",
			fields: fields{
				storage: successfulWithdrawal(),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "79927398713", Sum: 10},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balance{
				storage: tt.fields.storage,
			}
			if err := b.Withdraw(context.Background(), tt.args.userID, tt.args.withdrawal); (err != nil) != tt.wantErr {
				t.Errorf("Withdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func failingWithdrawal() *mocks.Storage {
	m := new(mocks.Storage)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("mock error"))
	return m
}

func successfulWithdrawal() *mocks.Storage {
	m := new(mocks.Storage)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(&model.Withdrawal{}, nil)
	return m
}
