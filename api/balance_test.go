package api

import (
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
	"testing"
)

func TestHandleBalance(t *testing.T) {
	type args struct {
		token   string
		balance balance.Balance
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
			name: "it returns 401 if user is not authorized",
			args: args{
				token:   "invalid token",
				balance: new(mocks.Balance),
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "it returns user's current and withdrawn balance",
			args: args{
				token:   token(100),
				balance: balanceForUserMock(100, 500, 42),
			},
			want: want{
				status:  http.StatusOK,
				headers: map[string]string{"Content-Type": "application/json"},
				body:    `{"current":500,"withdrawn":42}` + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(new(mocks.Auth), tt.args.balance, new(mocks.Order))
			require.NoError(t, err)

			r := a.routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest("GET", ts.URL+"/api/user/balance", nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.args.token))

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

func balanceForUserMock(userID uint64, current float64, withdrawn float64) *mocks.Balance {
	m := new(mocks.Balance)
	m.On("UserBalance", mock.Anything, userID).Return(
		&model.UserBalance{
			Current:   current,
			Withdrawn: withdrawn,
		},
		nil,
	)
	return m
}
