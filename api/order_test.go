package api

import (
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/soundrussian/go-practicum-diploma/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI_HandleOrder(t *testing.T) {
	type args struct {
		body    string
		token   string
		headers map[string]string
		order   order.Order
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
			name: "returns 415 if content type is not text/plain",
			args: args{
				order:   new(mocks.Order),
				token:   token(100),
				headers: map[string]string{"Content-Type": "application/json"},
				body:    "79927398713",
			},
			want: want{
				status: http.StatusUnsupportedMediaType,
			},
		},
		{
			name: "returns 422 if order id is invalid",
			args: args{
				token:   token(100),
				headers: map[string]string{"Content-Type": "text/plain"},
				order:   orderInvalidMock(),
				body:    "not an order id",
			},
			want: want{
				status: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "returns 409 if order was uploaded by another user",
			args: args{
				token:   token(100),
				headers: map[string]string{"Content-Type": "text/plain"},
				order:   orderUploadedByOtherUser(),
				body:    "79927398713",
			},
			want: want{
				status: http.StatusConflict,
			},
		},
		{
			name: "returns 200 if order was uploaded by this user",
			args: args{
				token:   token(100),
				headers: map[string]string{"Content-Type": "text/plain"},
				order:   orderUploadedByCurrentUser(),
				body:    "79927398713",
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "returns 202 if order is new",
			args: args{
				token:   token(100),
				headers: map[string]string{"Content-Type": "text/plain"},
				order:   orderSuccess(),
				body:    "79927398713",
			},
			want: want{
				status: http.StatusAccepted,
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

			req, err := http.NewRequest("POST", ts.URL+"/api/user/orders", strings.NewReader(tt.args.body))
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

func orderInvalidMock() *mocks.Order {
	m := new(mocks.Order)
	m.On("AcceptOrder", mock.Anything, mock.Anything, mock.Anything).Return(order.ErrOrderInvalid)
	return m
}

func orderUploadedByOtherUser() *mocks.Order {
	m := new(mocks.Order)
	m.On("AcceptOrder", mock.Anything, mock.Anything, mock.Anything).Return(order.ErrConflict)
	return m
}

func orderUploadedByCurrentUser() *mocks.Order {
	m := new(mocks.Order)
	m.On("AcceptOrder", mock.Anything, mock.Anything, mock.Anything).Return(order.ErrAlreadyAccepted)
	return m
}

func orderSuccess() *mocks.Order {
	m := new(mocks.Order)
	m.On("AcceptOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return m
}
