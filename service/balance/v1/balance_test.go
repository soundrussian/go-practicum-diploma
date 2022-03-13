package v1

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/service/balance"
	"github.com/soundrussian/go-practicum-diploma/storage"
	storageMock "github.com/soundrussian/go-practicum-diploma/storage/mock"
	"github.com/stretchr/testify/assert"
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
				storage: new(storageMock.Storage),
			},
			want:    &Balance{storage: new(storageMock.Storage)},
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
				storage: successfulStorage(decimal.NewFromInt(100), decimal.NewFromInt(50)),
			},
			args: args{
				userID: 100,
			},
			want:    &model.UserBalance{Current: decimal.NewFromInt(100), Withdrawn: decimal.NewFromInt(50)},
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

func failingStorage() *storageMock.Storage {
	m := new(storageMock.Storage)
	m.On("UserBalance", mock.Anything, mock.Anything).Return(nil, errors.New("mock error"))
	return m
}

func successfulStorage(current decimal.Decimal, withdrawn decimal.Decimal) *storageMock.Storage {
	m := new(storageMock.Storage)
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
		wantErr error
	}{
		{
			name: "returns error if sum is less than zero",
			fields: fields{
				storage: new(storageMock.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Sum: decimal.NewFromInt(-1)},
			},
			wantErr: balance.ErrInvalidSum,
		},
		{
			name: "returns error if sum is zero",
			fields: fields{
				storage: new(storageMock.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Sum: decimal.Zero},
			},
			wantErr: balance.ErrInvalidSum,
		},
		{
			name: "returns error if order number is missing",
			fields: fields{
				storage: new(storageMock.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Sum: decimal.NewFromInt(10)},
			},
			wantErr: balance.ErrInvalidOrder,
		},
		{
			name: "returns error if order number is not a number",
			fields: fields{
				storage: new(storageMock.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "not a number", Sum: decimal.NewFromInt(10)},
			},
			wantErr: balance.ErrInvalidOrder,
		},
		{
			name: "returns error if order checksum is invalid",
			fields: fields{
				storage: new(storageMock.Storage),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "7992739871", Sum: decimal.NewFromInt(10)},
			},
			wantErr: balance.ErrInvalidOrder,
		},
		{
			name: "returns error if storage reported error",
			fields: fields{
				storage: failingWithdrawal(),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "79927398713", Sum: decimal.NewFromInt(10)},
			},
			wantErr: balance.ErrInternalError,
		},
		{
			name: "returns ErrNotEnoughBalance if storage tells so",
			fields: fields{
				storage: notEnoughBalanceWithdrawal(),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "79927398713", Sum: decimal.NewFromInt(10)},
			},
			wantErr: balance.ErrNotEnoughBalance,
		},
		{
			name: "does not return error if storage reported success",
			fields: fields{
				storage: successfulWithdrawal(),
			},
			args: args{
				userID:     100,
				withdrawal: model.Withdrawal{Order: "79927398713", Sum: decimal.NewFromInt(10)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balance{
				storage: tt.fields.storage,
			}
			err := b.Withdraw(context.Background(), tt.args.userID, tt.args.withdrawal)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func failingWithdrawal() *storageMock.Storage {
	m := new(storageMock.Storage)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("mock error"))
	return m
}

func successfulWithdrawal() *storageMock.Storage {
	m := new(storageMock.Storage)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(&model.Withdrawal{}, nil)
	return m
}

func notEnoughBalanceWithdrawal() *storageMock.Storage {
	m := new(storageMock.Storage)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(nil, storage.ErrNotEnoughBalance)
	return m
}
