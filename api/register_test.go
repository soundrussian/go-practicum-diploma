package api

import (
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/auth/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	validRequest = `
		{
			"login": "john.doe@example.com",
			"password": "topsecret"
		}
	`
	invalidJSON = `not a json`
)

func TestHandleRegister(t *testing.T) {
	type args struct {
		body    string
		headers map[string]string
		auth    auth.Auth
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
			name: "returns 415 Unsupported Media Type if request does not have Content-Type: application/json",
			args: args{
				body: validRequest,
				auth: mock.SuccessfulRegistration{},
			},
			want: want{
				status: http.StatusUnsupportedMediaType,
			},
		},
		{
			name: "returns 400 Bad Request if request payload is not JSON",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: invalidJSON,
				auth: mock.SuccessfulRegistration{},
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "returns 400 Bad Request if Auth service could not register user",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: validRequest,
				auth: mock.FailedValidation{},
			},
			want: want{
				status: http.StatusBadRequest,
				body:   auth.ErrInvalidLogin.Error() + "\n",
			},
		},
		{
			name: "returns 409 Bad Request if user is already registered",
			args: args{
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: validRequest,
				auth: mock.DuplicateUser{},
			},
			want: want{
				status: http.StatusConflict,
				body:   auth.ErrUserAlreadyRegistered.Error() + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(tt.args.auth)
			require.NoError(t, err)

			r := a.routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest("POST", ts.URL+"/api/user/register", strings.NewReader(tt.args.body))
			require.NoError(t, err)

			for header, value := range tt.args.headers {
				req.Header.Set(header, value)
			}

			transport := http.Transport{}
			resp, err := transport.RoundTrip(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, resp.StatusCode)

			for wantHeader, wantHeaderValue := range tt.want.headers {
				assert.Equal(t, resp.Header.Get(wantHeader), wantHeaderValue)
			}

			if tt.want.body != "" {
				resBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tt.want.body, string(resBody))
			}
		})
	}
}
