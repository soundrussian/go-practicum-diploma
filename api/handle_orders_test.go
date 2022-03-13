package api

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/soundrussian/go-practicum-diploma/model"
	authMock "github.com/soundrussian/go-practicum-diploma/service/auth/mock"
	balanceMock "github.com/soundrussian/go-practicum-diploma/service/balance/mock"
	"github.com/soundrussian/go-practicum-diploma/service/order"
	orderMock "github.com/soundrussian/go-practicum-diploma/service/order/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
				order: new(orderMock.Order),
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
		{
			name: "returns 200 and orders list if there are some",
			args: args{
				token: token(100),
				order: ordersMock([]model.Order{
					{
						UserID:     100,
						Accrual:    decimal.Zero,
						OrderID:    "9278923470",
						Status:     model.OrderNew,
						UploadedAt: time.Date(2022, 3, 8, 9, 10, 11, 0, time.UTC),
					},
					{
						UserID:     100,
						Accrual:    decimal.Zero,
						OrderID:    "12345678903",
						Status:     model.OrderProcessing,
						UploadedAt: time.Date(2022, 3, 8, 9, 5, 11, 0, time.UTC),
					},
					{
						UserID:     100,
						Accrual:    decimal.Zero,
						OrderID:    "346436439",
						Status:     model.OrderInvalid,
						UploadedAt: time.Date(2022, 3, 7, 9, 5, 11, 0, time.UTC),
					},
					{
						UserID:     100,
						Accrual:    decimal.NewFromInt(500),
						OrderID:    "79927398713",
						Status:     model.OrderProcessed,
						UploadedAt: time.Date(2022, 3, 7, 9, 5, 11, 0, time.UTC),
					},
				}),
			},
			want: want{
				status:  http.StatusOK,
				headers: map[string]string{"Content-Type": "application/json"},
				body:    "testdata/orders.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(new(authMock.Auth), new(balanceMock.Balance), tt.args.order)
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
				body := []byte(tt.want.body)
				if strings.HasPrefix(tt.want.body, "testdata/") {
					body, err = ioutil.ReadFile(tt.want.body)
					require.NoError(t, err)
				}

				resBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, string(body), string(resBody))
			}

			assert.Equal(t, tt.want.status, resp.StatusCode)

			for wantHeader, wantHeaderValue := range tt.want.headers {
				assert.Equal(t, wantHeaderValue, resp.Header.Get(wantHeader))
			}

		})
	}
}

func noOrdersMock() *orderMock.Order {
	m := new(orderMock.Order)
	m.On("UserOrders", mock.Anything, mock.Anything).Return([]model.Order{}, nil)
	return m
}

func ordersMock(orders []model.Order) *orderMock.Order {
	m := new(orderMock.Order)
	m.On("UserOrders", mock.Anything, mock.Anything).Return(orders, nil)
	return m
}
