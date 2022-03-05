package api

import (
	"fmt"
	authMock "github.com/soundrussian/go-practicum-diploma/auth/mock"
	"github.com/soundrussian/go-practicum-diploma/balance"
	balanceMock "github.com/soundrussian/go-practicum-diploma/balance/mock"
	"github.com/stretchr/testify/assert"
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
				balance: balanceMock.BalanceMock{},
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "it returns user's current and withdrawn balance",
			args: args{
				token: authMock.Token(),
				balance: balanceMock.BalanceMock{
					ByUser: map[uint64]balance.UserBalance{
						authMock.UserID: {
							Current:   500,
							Withdrawn: 42,
						},
					},
				},
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
			a, err := New(authMock.Successful{}, tt.args.balance)
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
