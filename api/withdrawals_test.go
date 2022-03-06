package api

import (
	"errors"
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/balance"
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAPI_HandleWithdrawals(t *testing.T) {
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
			name: "returns 204 if there are no withdrawals",
			args: args{
				token:   token(100),
				balance: noWithdrawalsMock(),
			},
			want: want{
				status: http.StatusNoContent,
			},
		},
		{
			name: "returns 500 if there was error fetching withdrawals",
			args: args{
				token:   token(100),
				balance: failedWithdrawalMock(),
			},
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "returns 200 and list of withdrawals if no error",
			args: args{
				token: token(100),
				balance: successfulWithdrawalList([]model.Withdrawal{
					{
						Order:       "1",
						Sum:         100,
						ProcessedAt: time.Date(2022, 03, 06, 05, 04, 01, 0, time.UTC),
					},
					{
						Order:       "1",
						Sum:         25,
						ProcessedAt: time.Date(2022, 03, 06, 06, 04, 01, 0, time.UTC),
					},
				}),
			},
			want: want{
				status: http.StatusOK,
				body:   `[{"order":"1","sum":100,"processed_at":"2022-03-06T05:04:01Z"},{"order":"1","sum":25,"processed_at":"2022-03-06T06:04:01Z"}]` + "\n",
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

			req, err := http.NewRequest("GET", ts.URL+"/api/user/balance/withdrawals", strings.NewReader(tt.args.body))
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

func noWithdrawalsMock() *mocks.Balance {
	m := new(mocks.Balance)
	m.On("Withdrawals", mock.Anything, mock.Anything).Return([]model.Withdrawal{}, nil)
	return m
}

func failedWithdrawalMock() *mocks.Balance {
	m := new(mocks.Balance)
	m.On("Withdrawals", mock.Anything, mock.Anything).Return([]model.Withdrawal{}, errors.New("mock error"))
	return m
}

func successfulWithdrawalList(withdrawals []model.Withdrawal) *mocks.Balance {
	m := new(mocks.Balance)
	m.On("Withdrawals", mock.Anything, mock.Anything).Return(withdrawals, nil)
	return m
}
