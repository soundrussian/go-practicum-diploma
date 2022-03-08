package api

import (
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_HandleOrders(t *testing.T) {
	type args struct {
		token string
		order order.Order
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
				order: new(mocks.Order),
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "returns 204 if no orders for user",
			args: args{
				token: token(100),
				order: noOrdersMock(),
			},
			want: want{
				status: http.StatusNoContent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(new(mocks.Auth), new(mocks.Balance), tt.args.order)
			require.NoError(t, err)

			r := a.routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest("GET", ts.URL+"/api/user/orders", nil)
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

func noOrdersMock() *mocks.Order {
	m := new(mocks.Order)
	m.On("UserOrders", mock.Anything, mock.Anything).Return([]model.Order{}, nil)
	return m
}
