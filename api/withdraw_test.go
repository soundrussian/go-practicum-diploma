package api

import (
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/balance"
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI_HandleWithdraw(t *testing.T) {
	type args struct {
		body    string
		token   string
		balance balance.Balance
		headers map[string]string
	}
	type want struct {
		status  int
		headers map[string]string
		body    string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "returns 401 if user is not authorized",
			args: args{
				token:   "invalid token",
				balance: new(mocks.Balance),
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "returns 415 if content type is not application/json",
			args: args{
				token:   token(100),
				balance: new(mocks.Balance),
				body:    `not a json`,
			},
			want: want{
				status: http.StatusUnsupportedMediaType,
			},
		},
		{
			name: "returns 400 if body cannot be parsed",
			args: args{
				token:   token(100),
				balance: new(mocks.Balance),
				headers: map[string]string{"Content-Type": "application/json"},
				body:    `not a json`,
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "returns 422 if order is invalid",
			args: args{
				token:   token(100),
				balance: invalidOrderMock(),
				headers: map[string]string{"Content-Type": "application/json"},
				body:    `{"order": "not a number", "sum": 10}"`,
			},
			want: want{
				status: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "returns 422 if sum is negative",
			args: args{
				token:   token(100),
				balance: invalidSumMock(),
				headers: map[string]string{"Content-Type": "application/json"},
				body:    `{"order": "79927398713", "sum": -10}"`,
			},
			want: want{
				status: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "returns 402 if not enough balance for withdrawal",
			args: args{
				token:   token(100),
				balance: notEnoughBalanceMock(),
				headers: map[string]string{"Content-Type": "application/json"},
				body:    `{"order": "79927398713", "sum": 10}"`,
			},
			want: want{
				status: http.StatusPaymentRequired,
			},
		},
		{
			name: "returns 200 if everything is okay",
			args: args{
				token:   token(100),
				balance: successfulWithdrawal(),
				headers: map[string]string{"Content-Type": "application/json"},
				body:    `{"order": "79927398713", "sum": 10}"`,
			},
			want: want{
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(new(mocks.Auth), tt.args.balance)
			require.NoError(t, err)

			r := a.routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest("POST", ts.URL+"/api/user/balance/withdraw", strings.NewReader(tt.args.body))
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.args.token))

			for header, value := range tt.args.headers {
				req.Header.Set(header, value)
			}

			transport := http.Transport{}
			resp, err := transport.RoundTrip(req)
			require.NoError(t, err)

			if tt.want.body != "" {
				resBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tt.want.body, string(resBody))
			}

			assert.Equal(t, tt.want.status, resp.StatusCode)

			for wantHeader, wantHeaderValue := range tt.want.headers {
				assert.Equal(t, wantHeaderValue, resp.Header.Get(wantHeader))
			}

		})
	}
}

func notEnoughBalanceMock() *mocks.Balance {
	m := new(mocks.Balance)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(balance.ErrNotEnoughBalance)
	return m
}

func invalidSumMock() *mocks.Balance {
	m := new(mocks.Balance)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(balance.ErrInvalidSum)
	return m
}

func invalidOrderMock() *mocks.Balance {
	m := new(mocks.Balance)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(balance.ErrInvalidOrder)
	return m
}

func successfulWithdrawal() *mocks.Balance {
	m := new(mocks.Balance)
	m.On("Withdraw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return m
}
